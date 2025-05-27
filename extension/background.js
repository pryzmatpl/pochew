// Background service worker for Read It Later extension

const API_BASE_URL = 'http://localhost:8080/api';

// Initialize extension
chrome.runtime.onInstalled.addListener(() => {
  // Create context menu
  chrome.contextMenus.create({
    id: 'save-to-read-later',
    title: 'Save to Read It Later',
    contexts: ['page', 'selection', 'link']
  });

  console.log('Read It Later extension installed');
});

// Handle context menu clicks
chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  if (info.menuItemId === 'save-to-read-later') {
    try {
      await saveCurrentPage(tab, info);
    } catch (error) {
      console.error('Failed to save page:', error);
      showNotification('Failed to save page', 'error');
    }
  }
});

// Handle extension icon clicks
chrome.action.onClicked.addListener(async (tab) => {
  try {
    await saveCurrentPage(tab);
  } catch (error) {
    console.error('Failed to save page:', error);
    showNotification('Failed to save page', 'error');
  }
});

// Save current page
async function saveCurrentPage(tab, contextInfo = null) {
  // Check if user is authenticated
  const authToken = await getAuthToken();
  if (!authToken) {
    // Open login page
    chrome.tabs.create({ url: 'http://localhost:8080/login' });
    return;
  }

  // Extract content from the page
  const [result] = await chrome.scripting.executeScript({
    target: { tabId: tab.id },
    function: extractPageContent,
    args: [contextInfo]
  });

  if (!result || !result.result) {
    throw new Error('Failed to extract page content');
  }

  const pageData = result.result;

  // Save to backend
  const response = await fetch(`${API_BASE_URL}/articles`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authToken}`
    },
    body: JSON.stringify({
      title: pageData.title,
      url: pageData.url,
      content: pageData.content,
      summary: pageData.summary,
      tags: pageData.tags,
      metadata: pageData.metadata
    })
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const savedArticle = await response.json();
  
  // Show success notification
  showNotification(`Saved: ${pageData.title}`, 'success');
  
  // Update badge
  updateBadge(tab.id);
}

// Extract content from page (injected function)
function extractPageContent(contextInfo) {
  const result = {
    title: document.title,
    url: window.location.href,
    content: '',
    summary: '',
    tags: [],
    metadata: {}
  };

  // Extract main content
  const contentSelectors = [
    'article',
    '[role="main"]',
    '.content',
    '.post-content',
    '.entry-content',
    '.article-content',
    'main',
    '#content',
    '.main-content'
  ];

  let contentElement = null;
  for (const selector of contentSelectors) {
    contentElement = document.querySelector(selector);
    if (contentElement) break;
  }

  if (!contentElement) {
    contentElement = document.body;
  }

  // Clean and extract text content
  const clonedElement = contentElement.cloneNode(true);
  
  // Remove unwanted elements
  const unwantedSelectors = [
    'script', 'style', 'nav', 'header', 'footer', 
    '.advertisement', '.ads', '.social-share',
    '.comments', '.sidebar', '.menu'
  ];
  
  unwantedSelectors.forEach(selector => {
    const elements = clonedElement.querySelectorAll(selector);
    elements.forEach(el => el.remove());
  });

  result.content = clonedElement.innerText.trim();

  // Extract summary from meta description
  const metaDescription = document.querySelector('meta[name="description"]');
  if (metaDescription) {
    result.summary = metaDescription.getAttribute('content');
  }

  // Extract tags from meta keywords
  const metaKeywords = document.querySelector('meta[name="keywords"]');
  if (metaKeywords) {
    result.tags = metaKeywords.getAttribute('content')
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0);
  }

  // Extract additional metadata
  result.metadata = {
    author: getMetaContent('author'),
    publishedTime: getMetaContent('article:published_time') || getMetaContent('published_time'),
    siteName: getMetaContent('og:site_name'),
    type: getMetaContent('og:type'),
    image: getMetaContent('og:image'),
    readingTime: estimateReadingTime(result.content)
  };

  // If context info is provided (right-click), handle selection
  if (contextInfo && contextInfo.selectionText) {
    result.content = contextInfo.selectionText;
    result.summary = contextInfo.selectionText.substring(0, 200) + '...';
  }

  return result;
}

// Helper function to get meta content
function getMetaContent(name) {
  const meta = document.querySelector(`meta[name="${name}"], meta[property="${name}"]`);
  return meta ? meta.getAttribute('content') : null;
}

// Estimate reading time
function estimateReadingTime(text) {
  const wordsPerMinute = 200;
  const words = text.split(/\s+/).length;
  const minutes = Math.ceil(words / wordsPerMinute);
  return `${minutes} min read`;
}

// Get authentication token from storage
async function getAuthToken() {
  const result = await chrome.storage.local.get(['authToken']);
  return result.authToken;
}

// Set authentication token
async function setAuthToken(token) {
  await chrome.storage.local.set({ authToken: token });
}

// Clear authentication token
async function clearAuthToken() {
  await chrome.storage.local.remove(['authToken']);
}

// Show notification
function showNotification(message, type = 'info') {
  chrome.notifications.create({
    type: 'basic',
    iconUrl: 'icons/icon-48.png',
    title: 'Read It Later',
    message: message
  });
}

// Update badge
function updateBadge(tabId) {
  chrome.action.setBadgeText({
    text: '✓',
    tabId: tabId
  });
  
  chrome.action.setBadgeBackgroundColor({
    color: '#4CAF50',
    tabId: tabId
  });

  // Clear badge after 3 seconds
  setTimeout(() => {
    chrome.action.setBadgeText({
      text: '',
      tabId: tabId
    });
  }, 3000);
}

// Listen for messages from popup or content script
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  switch (request.action) {
    case 'saveCurrentPage':
      saveCurrentPage(sender.tab)
        .then(() => sendResponse({ success: true }))
        .catch(error => sendResponse({ success: false, error: error.message }));
      return true; // Keep message channel open for async response

    case 'getAuthToken':
      getAuthToken()
        .then(token => sendResponse({ token }))
        .catch(error => sendResponse({ error: error.message }));
      return true;

    case 'setAuthToken':
      setAuthToken(request.token)
        .then(() => sendResponse({ success: true }))
        .catch(error => sendResponse({ success: false, error: error.message }));
      return true;

    case 'clearAuthToken':
      clearAuthToken()
        .then(() => sendResponse({ success: true }))
        .catch(error => sendResponse({ success: false, error: error.message }));
      return true;

    default:
      sendResponse({ error: 'Unknown action' });
  }
});

// Handle web navigation to clear badges
chrome.webNavigation.onCompleted.addListener((details) => {
  if (details.frameId === 0) { // Main frame only
    chrome.action.setBadgeText({
      text: '',
      tabId: details.tabId
    });
  }
}); 
// Popup script for Read It Later extension

const API_BASE_URL = 'http://localhost:8080/api';

// DOM elements
const authSection = document.getElementById('auth-section');
const mainSection = document.getElementById('main-section');
const loginBtn = document.getElementById('login-btn');
const logoutBtn = document.getElementById('logout-btn');
const userEmail = document.getElementById('user-email');
const pageTitle = document.getElementById('page-title');
const pageUrl = document.getElementById('page-url');
const tagsInput = document.getElementById('tags-input');
const saveBtn = document.getElementById('save-btn');
const saveText = document.getElementById('save-text');
const status = document.getElementById('status');

// Current page data
let currentPageData = null;
let currentUser = null;

// Initialize popup
document.addEventListener('DOMContentLoaded', async () => {
  try {
    await initializePopup();
  } catch (error) {
    console.error('Failed to initialize popup:', error);
    showStatus('Failed to initialize extension', 'error');
  }
});

// Initialize popup state
async function initializePopup() {
  // Check authentication status
  const authToken = await getAuthToken();
  
  if (authToken) {
    // User is authenticated, get user info
    try {
      currentUser = await getCurrentUser(authToken);
      showMainSection();
      await loadCurrentPageInfo();
    } catch (error) {
      console.error('Failed to get user info:', error);
      // Token might be invalid, clear it
      await clearAuthToken();
      showAuthSection();
    }
  } else {
    showAuthSection();
  }
}

// Show authentication section
function showAuthSection() {
  authSection.classList.remove('hidden');
  mainSection.classList.add('hidden');
}

// Show main section
function showMainSection() {
  authSection.classList.add('hidden');
  mainSection.classList.remove('hidden');
  
  if (currentUser) {
    userEmail.textContent = currentUser.email;
  }
}

// Load current page information
async function loadCurrentPageInfo() {
  try {
    // Get current tab
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    
    if (!tab) {
      throw new Error('No active tab found');
    }

    // Extract page content
    const [result] = await chrome.scripting.executeScript({
      target: { tabId: tab.id },
      function: extractBasicPageInfo
    });

    if (result && result.result) {
      currentPageData = result.result;
      updatePageInfo();
    } else {
      throw new Error('Failed to extract page information');
    }
  } catch (error) {
    console.error('Failed to load page info:', error);
    pageTitle.textContent = 'Unable to load page information';
    pageUrl.textContent = 'Please try refreshing the page';
  }
}

// Extract basic page information (injected function)
function extractBasicPageInfo() {
  return {
    title: document.title || 'Untitled Page',
    url: window.location.href,
    domain: window.location.hostname
  };
}

// Update page info display
function updatePageInfo() {
  if (currentPageData) {
    pageTitle.textContent = currentPageData.title;
    pageUrl.textContent = currentPageData.url;
  }
}

// Get authentication token
async function getAuthToken() {
  return new Promise((resolve) => {
    chrome.runtime.sendMessage({ action: 'getAuthToken' }, (response) => {
      resolve(response.token);
    });
  });
}

// Set authentication token
async function setAuthToken(token) {
  return new Promise((resolve) => {
    chrome.runtime.sendMessage({ action: 'setAuthToken', token }, (response) => {
      resolve(response.success);
    });
  });
}

// Clear authentication token
async function clearAuthToken() {
  return new Promise((resolve) => {
    chrome.runtime.sendMessage({ action: 'clearAuthToken' }, (response) => {
      resolve(response.success);
    });
  });
}

// Get current user information
async function getCurrentUser(authToken) {
  const response = await fetch(`${API_BASE_URL}/user/profile`, {
    headers: {
      'Authorization': `Bearer ${authToken}`
    }
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return await response.json();
}

// Save current article
async function saveCurrentArticle() {
  if (!currentPageData) {
    showStatus('No page data available', 'error');
    return;
  }

  try {
    // Show loading state
    saveBtn.disabled = true;
    saveText.innerHTML = '<span class="loading-spinner"></span>Saving...';
    
    // Get auth token
    const authToken = await getAuthToken();
    if (!authToken) {
      throw new Error('Not authenticated');
    }

    // Get current tab for full content extraction
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    
    // Extract full content
    const [result] = await chrome.scripting.executeScript({
      target: { tabId: tab.id },
      function: extractFullPageContent
    });

    if (!result || !result.result) {
      throw new Error('Failed to extract page content');
    }

    const fullPageData = result.result;

    // Parse tags
    const tags = tagsInput.value
      .split(',')
      .map(tag => tag.trim())
      .filter(tag => tag.length > 0);

    // Prepare article data
    const articleData = {
      title: fullPageData.title,
      url: fullPageData.url,
      content: fullPageData.content,
      summary: fullPageData.summary,
      tags: tags,
      metadata: fullPageData.metadata
    };

    // Save to backend
    const response = await fetch(`${API_BASE_URL}/articles`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${authToken}`
      },
      body: JSON.stringify(articleData)
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.message || `HTTP error! status: ${response.status}`);
    }

    const savedArticle = await response.json();
    
    // Show success
    showStatus('Article saved successfully!', 'success');
    
    // Clear tags input
    tagsInput.value = '';
    
    // Update button text
    saveText.textContent = 'Saved ✓';
    
    // Reset button after delay
    setTimeout(() => {
      saveBtn.disabled = false;
      saveText.textContent = 'Save Article';
    }, 2000);

  } catch (error) {
    console.error('Failed to save article:', error);
    showStatus(`Failed to save: ${error.message}`, 'error');
    
    // Reset button
    saveBtn.disabled = false;
    saveText.textContent = 'Save Article';
  }
}

// Extract full page content (injected function)
function extractFullPageContent() {
  const result = {
    title: document.title,
    url: window.location.href,
    content: '',
    summary: '',
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
    '.comments', '.sidebar', '.menu', '.navigation'
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
  } else {
    // Fallback: use first paragraph or first 200 characters
    const firstParagraph = clonedElement.querySelector('p');
    if (firstParagraph) {
      result.summary = firstParagraph.innerText.substring(0, 200) + '...';
    } else {
      result.summary = result.content.substring(0, 200) + '...';
    }
  }

  // Extract metadata
  result.metadata = {
    author: getMetaContent('author'),
    publishedTime: getMetaContent('article:published_time') || getMetaContent('published_time'),
    siteName: getMetaContent('og:site_name'),
    type: getMetaContent('og:type'),
    image: getMetaContent('og:image'),
    readingTime: estimateReadingTime(result.content),
    wordCount: result.content.split(/\s+/).length
  };

  return result;

  // Helper functions
  function getMetaContent(name) {
    const meta = document.querySelector(`meta[name="${name}"], meta[property="${name}"]`);
    return meta ? meta.getAttribute('content') : null;
  }

  function estimateReadingTime(text) {
    const wordsPerMinute = 200;
    const words = text.split(/\s+/).length;
    const minutes = Math.ceil(words / wordsPerMinute);
    return `${minutes} min read`;
  }
}

// Show status message
function showStatus(message, type = 'info') {
  status.className = `status ${type}`;
  status.textContent = message;
  status.classList.remove('hidden');
  
  // Auto-hide after 5 seconds
  setTimeout(() => {
    status.classList.add('hidden');
  }, 5000);
}

// Event listeners
loginBtn.addEventListener('click', () => {
  // Open login page in new tab
  chrome.tabs.create({ url: 'http://localhost:8080/login' });
  window.close();
});

logoutBtn.addEventListener('click', async () => {
  try {
    await clearAuthToken();
    currentUser = null;
    showAuthSection();
    showStatus('Logged out successfully', 'success');
  } catch (error) {
    console.error('Failed to logout:', error);
    showStatus('Failed to logout', 'error');
  }
});

saveBtn.addEventListener('click', saveCurrentArticle);

// Handle Enter key in tags input
tagsInput.addEventListener('keypress', (e) => {
  if (e.key === 'Enter') {
    saveCurrentArticle();
  }
});

// Listen for authentication updates from background script
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'authUpdated') {
    initializePopup();
  }
}); 
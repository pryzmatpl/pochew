// Content script for Read It Later extension

// Add visual feedback when saving articles
let saveIndicator = null;

// Listen for messages from background script
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  switch (request.action) {
    case 'showSaveIndicator':
      showSaveIndicator(request.message);
      break;
    case 'hideSaveIndicator':
      hideSaveIndicator();
      break;
    default:
      break;
  }
});

// Show save indicator
function showSaveIndicator(message = 'Saving to Read It Later...') {
  // Remove existing indicator
  hideSaveIndicator();

  // Create indicator element
  saveIndicator = document.createElement('div');
  saveIndicator.id = 'read-it-later-indicator';
  saveIndicator.innerHTML = `
    <div style="
      position: fixed;
      top: 20px;
      right: 20px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      padding: 12px 20px;
      border-radius: 8px;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      font-weight: 500;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      z-index: 10000;
      display: flex;
      align-items: center;
      gap: 8px;
      animation: slideIn 0.3s ease-out;
    ">
      <div style="
        width: 16px;
        height: 16px;
        border: 2px solid #ffffff;
        border-radius: 50%;
        border-top-color: transparent;
        animation: spin 1s linear infinite;
      "></div>
      ${message}
    </div>
    <style>
      @keyframes slideIn {
        from {
          transform: translateX(100%);
          opacity: 0;
        }
        to {
          transform: translateX(0);
          opacity: 1;
        }
      }
      @keyframes spin {
        to { transform: rotate(360deg); }
      }
    </style>
  `;

  document.body.appendChild(saveIndicator);
}

// Hide save indicator
function hideSaveIndicator() {
  if (saveIndicator) {
    saveIndicator.remove();
    saveIndicator = null;
  }
}

// Show success indicator
function showSuccessIndicator(message = 'Article saved successfully!') {
  // Remove existing indicator
  hideSaveIndicator();

  // Create success indicator
  const successIndicator = document.createElement('div');
  successIndicator.innerHTML = `
    <div style="
      position: fixed;
      top: 20px;
      right: 20px;
      background: #10b981;
      color: white;
      padding: 12px 20px;
      border-radius: 8px;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      font-weight: 500;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      z-index: 10000;
      display: flex;
      align-items: center;
      gap: 8px;
      animation: slideIn 0.3s ease-out;
    ">
      <span style="font-size: 16px;">✓</span>
      ${message}
    </div>
    <style>
      @keyframes slideIn {
        from {
          transform: translateX(100%);
          opacity: 0;
        }
        to {
          transform: translateX(0);
          opacity: 1;
        }
      }
    </style>
  `;

  document.body.appendChild(successIndicator);

  // Auto-remove after 3 seconds
  setTimeout(() => {
    successIndicator.remove();
  }, 3000);
}

// Show error indicator
function showErrorIndicator(message = 'Failed to save article') {
  // Remove existing indicator
  hideSaveIndicator();

  // Create error indicator
  const errorIndicator = document.createElement('div');
  errorIndicator.innerHTML = `
    <div style="
      position: fixed;
      top: 20px;
      right: 20px;
      background: #dc2626;
      color: white;
      padding: 12px 20px;
      border-radius: 8px;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      font-weight: 500;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      z-index: 10000;
      display: flex;
      align-items: center;
      gap: 8px;
      animation: slideIn 0.3s ease-out;
    ">
      <span style="font-size: 16px;">✕</span>
      ${message}
    </div>
    <style>
      @keyframes slideIn {
        from {
          transform: translateX(100%);
          opacity: 0;
        }
        to {
          transform: translateX(0);
          opacity: 1;
        }
      }
    </style>
  `;

  document.body.appendChild(errorIndicator);

  // Auto-remove after 5 seconds
  setTimeout(() => {
    errorIndicator.remove();
  }, 5000);
}

// Add keyboard shortcut for quick save (Ctrl+Shift+S or Cmd+Shift+S)
document.addEventListener('keydown', (e) => {
  if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'S') {
    e.preventDefault();
    
    // Send message to background script to save current page
    chrome.runtime.sendMessage({ action: 'saveCurrentPage' }, (response) => {
      if (response.success) {
        showSuccessIndicator();
      } else {
        showErrorIndicator(response.error || 'Failed to save article');
      }
    });
  }
});

// Highlight selected text when right-clicking
let lastSelection = null;

document.addEventListener('mouseup', () => {
  const selection = window.getSelection();
  if (selection.toString().trim().length > 0) {
    lastSelection = {
      text: selection.toString(),
      range: selection.getRangeAt(0).cloneRange()
    };
  }
});

// Clean up when page unloads
window.addEventListener('beforeunload', () => {
  hideSaveIndicator();
}); 
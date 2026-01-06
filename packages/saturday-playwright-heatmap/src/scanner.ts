/**
 * Script to scan the DOM for interactable elements.
 * This runs in the browser context at the end of a test/step.
 */

export const scannerScript = `
(function() {
  function isVisible(el) {
    if (!el.offsetParent && el.tagName !== 'BODY') return false;
    const style = window.getComputedStyle(el);
    if (style.display === 'none' || style.visibility === 'hidden' || style.opacity === '0') return false;
    const rect = el.getBoundingClientRect();
    return rect.width > 0 && rect.height > 0;
  }

  function getSelector(el) {
    if (el.id) return '#' + el.id;
    let tagName = el.tagName.toLowerCase();
    
    // Attempt data-testid which is common in testing
    if (el.hasAttribute('data-testid')) {
        return '[' + 'data-testid="' + el.getAttribute('data-testid') + '"]';
    }

    if (el.className && typeof el.className === 'string') {
        const classes = el.className.split(/\\s+/).filter(Boolean).join('.');
        if (classes) return tagName + '.' + classes;
    }
    
    return tagName;
  }

  const interactableSelectors = [
    'button', 
    'a[href]', 
    'input:not([type="hidden"])', 
    'select', 
    'textarea', 
    '[role="button"]', 
    '[role="link"]',
    '[role="checkbox"]',
    '[role="menuitem"]',
    '[role="tab"]',
    '[contenteditable="true"]'
  ];

  const elements = document.querySelectorAll(interactableSelectors.join(','));
  const results = [];

  elements.forEach(el => {
    if (isVisible(el) && !el.disabled) {
        const rect = el.getBoundingClientRect();
        results.push({
            selector: getSelector(el),
            tagName: el.tagName.toLowerCase(),
            text: el.innerText ? el.innerText.substring(0, 50) : '',
            rect: {
                x: rect.x + window.scrollX,
                y: rect.y + window.scrollY,
                width: rect.width,
                height: rect.height
            }
        });
    }
  });

  return results;
})();
`;

/**
 * Simple script to track basic interactions and store them on the window object.
 * This runs in the browser context.
 */

declare global {
  interface Window {
    __SATURDAY_HEATMAP_EVENTS__?: Array<{
      type: string;
      x: number;
      y: number;
      selector: string;
      timestamp: number;
    }>;
  }
}

export function generateSelector(el: Element): string {
  if (el.id) {
    return `#${el.id}`;
  }
  const tagName = el.tagName.toLowerCase();
  if (tagName === 'body' || tagName === 'html') {
    return tagName;
  }
  
  // Try using class names
  // Use 'Array.from' for broader compatibility if needed, but modern browsers support classList
  if (el.className && typeof el.className === 'string') {
     const classes = el.className.split(/\s+/).filter(c => c).join('.');
     if (classes) {
         return `${tagName}.${classes}`;
     }
  }

  // Fallback to simpler structural path if no ID or helpful classes
  let path = tagName;
  if (el.parentElement) {
    const siblings = Array.from(el.parentElement.children).filter(e => e.tagName === el.tagName);
    if (siblings.length > 1) {
        const index = siblings.indexOf(el) + 1;
        path += `:nth-of-type(${index})`;
    }
  }
  return path;
}

export const trackerScript = `
(function() {
  if (window.__SATURDAY_HEATMAP_INITIALIZED__) return;
  window.__SATURDAY_HEATMAP_INITIALIZED__ = true;
  window.__SATURDAY_HEATMAP_EVENTS__ = [];

  function getSelector(el) {
    if (el.id) return '#' + el.id;
    let tagName = el.tagName.toLowerCase();
    if (tagName === 'body' || tagName === 'html') return tagName;
    
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

  function recordEvent(e) {
    try {
        const selector = getSelector(e.target);
        window.__SATURDAY_HEATMAP_EVENTS__.push({
            type: e.type,
            x: e.clientX,
            y: e.clientY,
            selector: selector,
            timestamp: Date.now(),
            href: window.location.href
        });
    } catch(err) {
        console.error('Heatmap tracker error:', err);
    }
  }

  ['click', 'input', 'change'].forEach(type => {
    window.addEventListener(type, recordEvent, { capture: true, passive: true });
  });
})();
`;

/**
 * Helper utility functions
 */

/**
 * Format date to readable string
 */
function getLocale() {
  return typeof window !== 'undefined' && window.SiteI18n?.language === 'en' ? 'en-US' : 'vi-VN';
}

function formatDate(date, format = 'dd/mm/yyyy') {
  if (!date) return '';
  const d = new Date(date);
  const locale = getLocale();
  if (Number.isNaN(d.getTime())) return d.toLocaleString(locale);

  if (format === 'dd/mm/yyyy') {
    return new Intl.DateTimeFormat(locale, { day: '2-digit', month: '2-digit', year: 'numeric' }).format(d);
  }
  if (format === 'dd/mm/yyyy hh:mm') {
    return new Intl.DateTimeFormat(locale, {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(d);
  }
  return d.toLocaleString(locale);
}

/**
 * Format currency
 */
function formatCurrency(amount) {
  return new Intl.NumberFormat(getLocale(), {
    style: 'currency',
    currency: 'VND',
  }).format(amount);
}

/**
 * Format number with thousand separators
 */
function formatNumber(num) {
  return new Intl.NumberFormat(getLocale()).format(num);
}

/**
 * Debounce function
 */
function debounce(fn, delay) {
  let timeoutId;
  return function (...args) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn.apply(this, args), delay);
  };
}

/**
 * Throttle function
 */
function throttle(fn, limit) {
  let inThrottle;
  return function (...args) {
    if (!inThrottle) {
      fn.apply(this, args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

/**
 * Deep clone object
 */
function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') return obj;
  if (obj instanceof Date) return new Date(obj.getTime());
  if (obj instanceof Array) return obj.map((item) => deepClone(item));
  if (obj instanceof Object) {
    const cloned = {};
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        cloned[key] = deepClone(obj[key]);
      }
    }
    return cloned;
  }
}

/**
 * Check if email is valid
 */
function isValidEmail(email) {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return re.test(email);
}

/**
 * Check if phone is valid (Vietnam format)
 */
function isValidPhone(phone) {
  const re = /^(0|84)[0-9]{9}$/;
  return re.test(phone.replace(/\s/g, ''));
}

/**
 * Get query parameter from URL
 */
function getQueryParam(param) {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(param);
}

/**
 * Show toast notification
 */
function showToast(message, type = 'info', duration = 3000) {
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.textContent = message;
  toast.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    padding: 15px 20px;
    border-radius: 4px;
    background-color: ${type === 'success' ? '#4CAF50' : type === 'error' ? '#f44336' : type === 'warning' ? '#ff9800' : '#2196F3'};
    color: white;
    font-size: 14px;
    z-index: 10000;
    animation: slideIn 0.3s ease-in-out;
  `;

  document.body.appendChild(toast);
  setTimeout(() => toast.remove(), duration);
}

/**
 * Delay/sleep function
 */
function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export {
  formatDate,
  formatCurrency,
  formatNumber,
  debounce,
  throttle,
  deepClone,
  isValidEmail,
  isValidPhone,
  getQueryParam,
  showToast,
  sleep,
};

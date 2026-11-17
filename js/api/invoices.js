/**
 * Invoice API endpoints
 */

import apiClient from './client.js';
import { API_CONFIG } from '../utils/constants.js';

const invoiceAPI = {
  /**
   * Get all invoices
   */
  async getInvoices() {
    return apiClient.get('/invoices');
  },

  /**
   * Get invoice by ID
   */
  async getInvoiceByID(id) {
    return apiClient.get(`/invoices/${id}`);
  },

  /**
   * Create invoice (Admin/Manager only)
   */
  async createInvoice(data) {
    return apiClient.post('/invoices', data);
  },

  /**
   * Update invoice (Admin/Manager only)
   */
  async updateInvoice(id, data) {
    return apiClient.put(`/invoices/${id}`, data);
  },

  /**
   * Record payment (Staff/Manager only)
   */
  async recordPayment(invoiceID, data) {
    return apiClient.post(`/invoices/${invoiceID}/payments`, data);
  },

  /**
   * Initiate VNPay online payment
   */
  async initiateVNPay(invoiceID) {
    return apiClient.post(`/invoices/${invoiceID}/payments/vnpay`);
  },

  /**
   * Get payment transaction status
   */
  async getPaymentStatus(paymentID) {
    return apiClient.get(`/payments/${paymentID}`);
  },

  /**
   * Verify the signed VNPay return callback without requiring an active session.
   */
  async verifyVNPayReturn(queryString) {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), API_CONFIG.TIMEOUT);
    try {
      const response = await fetch(
        `${apiClient.baseURL}/payments/vnpay/return?${queryString}`,
        { headers: { Accept: 'application/json' }, signal: controller.signal }
      );
      const result = await response.json();
      if (!response.ok) {
        throw new Error(result.error || `HTTP ${response.status}`);
      }
      return result;
    } catch (error) {
      if (error.name === 'AbortError') {
        throw new Error('Timed out while verifying the VNPay payment');
      }
      throw error;
    } finally {
      clearTimeout(timeoutId);
    }
  },

  /**
   * Get overdue invoices
   */
  async getOverdueInvoices() {
    return apiClient.get('/invoices/overdue');
  },

  /**
   * Get financial report
   */
  async getFinancialReport(data) {
    return apiClient.post('/reports/financial', data);
  },
};

export default invoiceAPI;

// User Management for Media Vault Admin

class UserManager {
  constructor() {
    this.config = window.APP_CONFIG || {};
    this.endpoints = {
      users: `${this.config.API_BASE_URL}/api/admin/users`,
      roles: `${this.config.API_BASE_URL}/api/admin/roles`
    };
    this.currentPage = 1;
    this.pageSize = 10;
    this.totalUsers = 0;
    this.users = [];
    this.roles = [];
    
    // Initialize if on users page
    if (document.getElementById('users-page')) {
      this.init();
    }
  }

  /**
   * Initialize the user management interface
   */
  async init() {
    this.bindEvents();
    await this.loadRoles();
    await this.loadUsers();
    this.renderUserTable();
  }

  /**
   * Bind event listeners
   */
  bindEvents() {
    // User search
    const searchInput = document.getElementById('user-search');
    if (searchInput) {
      searchInput.addEventListener('input', this.debounce(() => {
        this.currentPage = 1;
        this.loadUsers();
      }, 300));
    }

    // Add user modal
    const addUserBtn = document.getElementById('add-user-btn');
    if (addUserBtn) {
      addUserBtn.addEventListener('click', () => this.showAddUserModal());
    }

    // Save user form
    const userForm = document.getElementById('user-form');
    if (userForm) {
      userForm.addEventListener('submit', (e) => {
        e.preventDefault();
        this.saveUser();
      });
    }

    // Role filter
    const roleFilter = document.getElementById('role-filter');
    if (roleFilter) {
      roleFilter.addEventListener('change', () => {
        this.currentPage = 1;
        this.loadUsers();
      });
    }

    // Pagination
    document.addEventListener('click', (e) => {
      if (e.target.closest('.page-link')) {
        e.preventDefault();
        const page = parseInt(e.target.closest('.page-link').dataset.page);
        if (page && !isNaN(page)) {
          this.currentPage = page;
          this.loadUsers();
        }
      }
    });
  }

  /**
   * Load users from the API
   */
  async loadUsers() {
    try {
      this.showLoading(true);
      
      const search = document.getElementById('user-search')?.value || '';
      const role = document.getElementById('role-filter')?.value || '';
      
      let url = `${this.endpoints.users}?page=${this.currentPage}&pageSize=${this.pageSize}`;
      if (search) url += `&search=${encodeURIComponent(search)}`;
      if (role) url += `&role=${encodeURIComponent(role)}`;
      
      const response = await fetch(url, {
        headers: this.getAuthHeaders()
      });
      
      if (!response.ok) {
        throw new Error('Failed to load users');
      }
      
      const data = await response.json();
      this.users = data.users || [];
      this.totalUsers = data.total || 0;
      
      this.renderUserTable();
      this.renderPagination();
      
    } catch (error) {
      console.error('Error loading users:', error);
      this.showError('Failed to load users. Please try again.');
    } finally {
      this.showLoading(false);
    }
  }

  /**
   * Load available roles from the API
   */
  async loadRoles() {
    try {
      const response = await fetch(this.endpoints.roles, {
        headers: this.getAuthHeaders()
      });
      
      if (!response.ok) {
        throw new Error('Failed to load roles');
      }
      
      this.roles = await response.json();
      this.renderRoleFilter();
      
    } catch (error) {
      console.error('Error loading roles:', error);
      this.showError('Failed to load roles. Some features may be limited.');
    }
  }

  /**
   * Render the users table
   */
  renderUserTable() {
    const tbody = document.querySelector('#users-table tbody');
    if (!tbody) return;
    
    tbody.innerHTML = '';
    
    if (this.users.length === 0) {
      const tr = document.createElement('tr');
      tr.innerHTML = `
        <td colspan="6" class="text-center py-4 text-muted">
          No users found
        </td>
      `;
      tbody.appendChild(tr);
      return;
    }
    
    this.users.forEach(user => {
      const tr = document.createElement('tr');
      tr.className = 'align-middle';
      
      const statusClass = user.isActive ? 'success' : 'secondary';
      const statusText = user.isActive ? 'Active' : 'Inactive';
      
      tr.innerHTML = `
        <td>
          <div class="d-flex align-items-center">
            <div class="avatar avatar-sm me-3">
              <span class="avatar-text rounded-circle bg-light text-dark">
                ${user.name?.charAt(0).toUpperCase() || 'U'}
              </span>
            </div>
            <div>
              <div class="fw-semibold">${user.name || 'N/A'}</div>
              <div class="text-muted small">${user.email || ''}</div>
            </div>
          </div>
        </td>
        <td>${user.username || 'N/A'}</td>
        <td>
          ${user.roles?.map(role => 
            `<span class="badge bg-primary me-1">${role}</span>`
          ).join('') || '<span class="text-muted">No roles</span>'}
        </td>
        <td>
          <span class="badge bg-${statusClass}">${statusText}</span>
        </td>
        <td>${new Date(user.createdAt).toLocaleDateString()}</td>
        <td class="text-end">
          <div class="dropdown">
            <button class="btn btn-sm btn-link text-dark" type="button" data-bs-toggle="dropdown" aria-expanded="false">
              <i class="bi bi-three-dots-vertical"></i>
            </button>
            <ul class="dropdown-menu dropdown-menu-end">
              <li>
                <a class="dropdown-item" href="#" data-action="edit" data-user-id="${user.id}">
                  <i class="bi bi-pencil me-2"></i>Edit
                </a>
              </li>
              <li>
                <a class="dropdown-item" href="#" data-action="reset-password" data-user-id="${user.id}">
                  <i class="bi bi-key me-2"></i>Reset Password
                </a>
              </li>
              <li><hr class="dropdown-divider"></li>
              <li>
                <a class="dropdown-item text-danger" href="#" data-action="delete" data-user-id="${user.id}">
                  <i class="bi bi-trash me-2"></i>Delete
                </a>
              </li>
            </ul>
          </div>
        </td>
      `;
      
      tbody.appendChild(tr);
    });
    
    // Add event listeners to action buttons
    document.querySelectorAll('[data-action]').forEach(button => {
      button.addEventListener('click', (e) => {
        e.preventDefault();
        const action = button.dataset.action;
        const userId = button.dataset.userId;
        
        switch (action) {
          case 'edit':
            this.editUser(userId);
            break;
          case 'delete':
            this.confirmDeleteUser(userId);
            break;
          case 'reset-password':
            this.resetPassword(userId);
            break;
        }
      });
    });
  }

  /**
   * Render the role filter dropdown
   */
  renderRoleFilter() {
    const roleFilter = document.getElementById('role-filter');
    if (!roleFilter) return;
    
    // Clear existing options except the first one
    while (roleFilter.options.length > 1) {
      roleFilter.remove(1);
    }
    
    // Add role options
    this.roles.forEach(role => {
      const option = document.createElement('option');
      option.value = role.id;
      option.textContent = role.name;
      roleFilter.appendChild(option);
    });
  }

  /**
   * Render pagination controls
   */
  renderPagination() {
    const pagination = document.getElementById('pagination');
    if (!pagination) return;
    
    const totalPages = Math.ceil(this.totalUsers / this.pageSize);
    
    if (totalPages <= 1) {
      pagination.innerHTML = '';
      return;
    }
    
    let html = `
      <nav aria-label="User pagination">
        <ul class="pagination mb-0">
          <li class="page-item ${this.currentPage === 1 ? 'disabled' : ''}">
            <a class="page-link" href="#" data-page="${this.currentPage - 1}" aria-label="Previous">
              <span aria-hidden="true">&laquo;</span>
            </a>
          </li>
    `;
    
    // Show first page, current page, and pages around it
    const pagesToShow = 3; // Number of pages to show around current page
    let startPage = Math.max(1, this.currentPage - pagesToShow);
    let endPage = Math.min(totalPages, this.currentPage + pagesToShow);
    
    // Adjust if we're near the start or end
    if (this.currentPage <= pagesToShow) {
      endPage = Math.min(2 * pagesToShow + 1, totalPages);
    }
    
    if (totalPages - this.currentPage < pagesToShow) {
      startPage = Math.max(1, totalPages - 2 * pagesToShow);
    }
    
    // First page
    if (startPage > 1) {
      html += `
        <li class="page-item ${1 === this.currentPage ? 'active' : ''}">
          <a class="page-link" href="#" data-page="1">1</a>
        </li>
      `;
      
      if (startPage > 2) {
        html += '<li class="page-item disabled"><span class="page-link">...</span></li>';
      }
    }
    
    // Page numbers
    for (let i = startPage; i <= endPage; i++) {
      html += `
        <li class="page-item ${i === this.currentPage ? 'active' : ''}">
          <a class="page-link" href="#" data-page="${i}">${i}</a>
        </li>
      `;
    }
    
    // Last page
    if (endPage < totalPages) {
      if (endPage < totalPages - 1) {
        html += '<li class="page-item disabled"><span class="page-link">...</span></li>';
      }
      
      html += `
        <li class="page-item ${totalPages === this.currentPage ? 'active' : ''}">
          <a class="page-link" href="#" data-page="${totalPages}">${totalPages}</a>
        </li>
      `;
    }
    
    // Next button
    html += `
          <li class="page-item ${this.currentPage === totalPages ? 'disabled' : ''}">
            <a class="page-link" href="#" data-page="${this.currentPage + 1}" aria-label="Next">
              <span aria-hidden="true">&raquo;</span>
            </a>
          </li>
        </ul>
      </nav>
    `;
    
    pagination.innerHTML = html;
  }

  /**
   * Show the add user modal
   */
  showAddUserModal(user = null) {
    const modal = new bootstrap.Modal(document.getElementById('userModal'));
    const form = document.getElementById('user-form');
    const modalTitle = document.getElementById('userModalLabel');
    const userIdField = document.getElementById('user-id');
    const isEdit = !!user;
    
    // Reset form
    form.reset();
    
    // Set modal title and button text
    if (isEdit) {
      modalTitle.textContent = 'Edit User';
      document.getElementById('save-user-btn').textContent = 'Update User';
      
      // Populate form fields
      userIdField.value = user.id;
      document.getElementById('name').value = user.name || '';
      document.getElementById('email').value = user.email || '';
      document.getElementById('username').value = user.username || '';
      document.getElementById('isActive').checked = user.isActive !== false;
      
      // Set roles
      const roleSelect = document.getElementById('roles');
      Array.from(roleSelect.options).forEach(option => {
        option.selected = user.roles?.includes(option.value) || false;
      });
      
    } else {
      modalTitle.textContent = 'Add New User';
      document.getElementById('save-user-btn').textContent = 'Add User';
      userIdField.value = '';
    }
    
    modal.show();
  }

  /**
   * Edit a user
   */
  async editUser(userId) {
    try {
      this.showLoading(true);
      
      const response = await fetch(`${this.endpoints.users}/${userId}`, {
        headers: this.getAuthHeaders()
      });
      
      if (!response.ok) {
        throw new Error('Failed to load user data');
      }
      
      const user = await response.json();
      this.showAddUserModal(user);
      
    } catch (error) {
      console.error('Error loading user:', error);
      this.showError('Failed to load user data. Please try again.');
    } finally {
      this.showLoading(false);
    }
  }

  /**
   * Save user (create or update)
   */
  async saveUser() {
    const form = document.getElementById('user-form');
    if (!form) return;
    
    const userId = document.getElementById('user-id').value;
    const isEdit = !!userId;
    
    const formData = new FormData(form);
    const userData = {
      name: formData.get('name'),
      email: formData.get('email'),
      username: formData.get('username'),
      isActive: formData.get('isActive') === 'on',
      roles: Array.from(document.getElementById('roles').selectedOptions).map(opt => opt.value)
    };
    
    // Add password only for new users
    if (!isEdit) {
      const password = formData.get('password');
      const confirmPassword = formData.get('confirmPassword');
      
      if (!password) {
        this.showError('Password is required');
        return;
      }
      
      if (password !== confirmPassword) {
        this.showError('Passwords do not match');
        return;
      }
      
      userData.password = password;
    }
    
    try {
      this.showLoading(true);
      
      const url = isEdit 
        ? `${this.endpoints.users}/${userId}`
        : this.endpoints.users;
      
      const method = isEdit ? 'PUT' : 'POST';
      
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          ...this.getAuthHeaders()
        },
        body: JSON.stringify(userData)
      });
      
      if (!response.ok) {
        const error = await response.json().catch(() => ({}));
        throw new Error(error.message || 'Failed to save user');
      }
      
      // Close modal and refresh user list
      const modal = bootstrap.Modal.getInstance(document.getElementById('userModal'));
      if (modal) modal.hide();
      
      this.showSuccess(`User ${isEdit ? 'updated' : 'created'} successfully`);
      this.loadUsers();
      
    } catch (error) {
      console.error('Error saving user:', error);
      this.showError(error.message || 'Failed to save user. Please try again.');
    } finally {
      this.showLoading(false);
    }
  }

  /**
   * Confirm user deletion
   */
  confirmDeleteUser(userId) {
    const user = this.users.find(u => u.id === userId);
    if (!user) return;
    
    if (confirm(`Are you sure you want to delete the user "${user.name || user.username}"? This action cannot be undone.`)) {
      this.deleteUser(userId);
    }
  }

  /**
   * Delete a user
   */
  async deleteUser(userId) {
    try {
      this.showLoading(true);
      
      const response = await fetch(`${this.endpoints.users}/${userId}`, {
        method: 'DELETE',
        headers: this.getAuthHeaders()
      });
      
      if (!response.ok) {
        throw new Error('Failed to delete user');
      }
      
      this.showSuccess('User deleted successfully');
      this.loadUsers();
      
    } catch (error) {
      console.error('Error deleting user:', error);
      this.showError('Failed to delete user. Please try again.');
    } finally {
      this.showLoading(false);
    }
  }

  /**
   * Reset user password
   */
  async resetPassword(userId) {
    const user = this.users.find(u => u.id === userId);
    if (!user) return;
    
    const newPassword = prompt(`Enter new password for ${user.name || user.username}:`);
    if (!newPassword) return;
    
    try {
      this.showLoading(true);
      
      const response = await fetch(`${this.endpoints.users}/${userId}/reset-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...this.getAuthHeaders()
        },
        body: JSON.stringify({ newPassword })
      });
      
      if (!response.ok) {
        throw new Error('Failed to reset password');
      }
      
      this.showSuccess('Password reset successfully');
      
    } catch (error) {
      console.error('Error resetting password:', error);
      this.showError('Failed to reset password. Please try again.');
    } finally {
      this.showLoading(false);
    }
  }

  /**
   * Show loading state
   */
  showLoading(isLoading) {
    const buttons = document.querySelectorAll('button[type="submit"]');
    buttons.forEach(button => {
      button.disabled = isLoading;
      const spinner = button.querySelector('.spinner-border');
      if (spinner) {
        spinner.style.display = isLoading ? 'inline-block' : 'none';
      }
    });
  }

  /**
   * Show success message
   */
  showSuccess(message) {
    // You can implement a toast notification system here
    alert(`Success: ${message}`);
  }

  /**
   * Show error message
   */
  showError(message) {
    // You can implement a toast notification system here
    alert(`Error: ${message}`);
  }

  /**
   * Get authentication headers
   */
  getAuthHeaders() {
    const token = localStorage.getItem('auth_token');
    return token ? { 'Authorization': `Bearer ${token}` } : {};
  }

  /**
   * Debounce function to limit the rate of function calls
   */
  debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  }
}

// Initialize the user manager when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', () => {
  window.userManager = new UserManager();
});
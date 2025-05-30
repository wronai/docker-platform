// Media Vault Admin - Main JavaScript

class MediaVaultAdmin {
  constructor() {
    this.config = window.APP_CONFIG || {};
    this.endpoints = {
      stats: `${this.config.API_BASE_URL}/api/admin/stats`,
      activity: `${this.config.API_BASE_URL}/api/admin/activity`,
      system: `${this.config.API_BASE_URL}/api/admin/system`,
    };
    this.init();
  }

  /**
   * Initialize the admin interface
   */
  init() {
    this.bindEvents();
    this.loadDashboardData();
    this.setupCharts();
  }

  /**
   * Bind event listeners
   */
  bindEvents() {
    // Sidebar toggle
    document.getElementById('sidebarCollapse')?.addEventListener('click', () => {
      document.getElementById('sidebar')?.classList.toggle('active');
      document.getElementById('content')?.classList.toggle('active');
    });

    // Logout button
    document.getElementById('logout')?.addEventListener('click', (e) => {
      e.preventDefault();
      this.logout();
    });

    // Refresh stats button
    document.getElementById('refresh-stats')?.addEventListener('click', () => {
      this.loadDashboardData();
    });
  }

  /**
   * Load dashboard data
   */
  async loadDashboardData() {
    try {
      // Show loading state
      this.setLoadingState(true);
      
      // Fetch stats
      const statsResponse = await fetch(this.endpoints.stats, {
        headers: this.getAuthHeaders()
      });
      
      if (!statsResponse.ok) {
        throw new Error('Failed to load dashboard stats');
      }
      
      const stats = await statsResponse.json();
      this.updateStatsUI(stats);
      
      // Fetch recent activity
      const activityResponse = await fetch(this.endpoints.activity, {
        headers: this.getAuthHeaders()
      });
      
      if (!activityResponse.ok) {
        throw new Error('Failed to load recent activity');
      }
      
      const activities = await activityResponse.json();
      this.updateActivityUI(activities);
      
      // Fetch system status
      const systemResponse = await fetch(this.endpoints.system, {
        headers: this.getAuthHeaders()
      });
      
      if (!systemResponse.ok) {
        throw new Error('Failed to load system status');
      }
      
      const systemStatus = await systemResponse.json();
      this.updateSystemStatusUI(systemStatus);
      
    } catch (error) {
      console.error('Error loading dashboard data:', error);
      this.showError('Failed to load dashboard data. Please try again.');
    } finally {
      this.setLoadingState(false);
    }
  }

  /**
   * Update stats cards UI
   */
  updateStatsUI(stats) {
    // Total files
    if (document.getElementById('total-files')) {
      document.getElementById('total-files').textContent = stats.totalFiles?.toLocaleString() || '0';
    }
    
    // Active users
    if (document.getElementById('active-users')) {
      document.getElementById('active-users').textContent = stats.activeUsers?.toLocaleString() || '0';
    }
    
    // Storage used
    if (document.getElementById('storage-used')) {
      const storageGB = (stats.storageUsed / 1024 / 1024 / 1024).toFixed(2);
      document.getElementById('storage-used').textContent = `${storageGB} GB`;
    }
    
    // Alerts
    if (document.getElementById('alerts-count')) {
      document.getElementById('alerts-count').textContent = stats.alertsCount?.toLocaleString() || '0';
    }
  }

  /**
   * Update activity feed UI
   */
  updateActivityUI(activities = []) {
    const activityList = document.getElementById('activity-list');
    if (!activityList) return;
    
    // Clear existing content
    activityList.innerHTML = '';
    
    // Add new activities
    activities.slice(0, 5).forEach(activity => {
      const activityItem = document.createElement('tr');
      activityItem.className = 'fade-in';
      
      const timeAgo = this.formatTimeAgo(activity.timestamp);
      
      activityItem.innerHTML = `
        <td class="align-middle">
          <div class="d-flex align-items-center">
            <div class="avatar avatar-sm me-2">
              <span class="avatar-text rounded-circle bg-light text-dark">
                ${activity.user?.charAt(0).toUpperCase() || 'U'}
              </span>
            </div>
            <div>
              <div class="fw-semibold">${activity.user || 'System'}</div>
              <small class="text-muted">${timeAgo}</small>
            </div>
          </div>
        </td>
        <td class="align-middle">
          <span class="badge bg-light text-dark">${activity.action}</span>
        </td>
        <td class="align-middle text-truncate" style="max-width: 200px;">
          ${activity.target || 'N/A'}
        </td>
        <td class="align-middle text-end">
          <small class="text-muted">${new Date(activity.timestamp).toLocaleTimeString()}</small>
        </td>
      `;
      
      activityList.appendChild(activityItem);
    });
  }

  /**
   * Update system status UI
   */
  updateSystemStatusUI(status) {
    // CPU Usage
    if (status.cpu) {
      const cpuUsage = Math.round(status.cpu.usage * 100);
      document.getElementById('cpu-usage').textContent = `${cpuUsage}%`;
      const cpuProgress = document.getElementById('cpu-progress');
      if (cpuProgress) {
        cpuProgress.style.width = `${cpuUsage}%`;
        cpuProgress.className = `progress-bar ${this.getUsageColorClass(cpuUsage)}`;
      }
    }
    
    // Memory Usage
    if (status.memory) {
      const memoryUsage = Math.round((status.memory.used / status.memory.total) * 100);
      document.getElementById('memory-usage').textContent = `${memoryUsage}%`;
      const memoryProgress = document.getElementById('memory-progress');
      if (memoryProgress) {
        memoryProgress.style.width = `${memoryUsage}%`;
        memoryProgress.className = `progress-bar ${this.getUsageColorClass(memoryUsage)}`;
      }
    }
    
    // Storage
    if (status.storage) {
      const storagePercent = Math.round((status.storage.used / status.storage.total) * 100);
      document.getElementById('storage-percent').textContent = `${storagePercent}%`;
      const storageProgress = document.getElementById('storage-progress');
      if (storageProgress) {
        storageProgress.style.width = `${storagePercent}%`;
        storageProgress.className = `progress-bar ${this.getUsageColorClass(storagePercent)}`;
      }
      
      // Update storage details
      const usedGB = (status.storage.used / 1024 / 1024 / 1024).toFixed(1);
      const totalGB = (status.storage.total / 1024 / 1024 / 1024).toFixed(1);
      document.getElementById('storage-details').textContent = `${usedGB} GB of ${totalGB} GB used`;
    }
  }

  /**
   * Setup charts
   */
  setupCharts() {
    // Check if Chart.js is available
    if (typeof Chart === 'undefined') {
      console.warn('Chart.js is not loaded. Charts will not be rendered.');
      return;
    }
    
    // Example chart setup - you can add more charts as needed
    const ctx = document.getElementById('usage-chart');
    if (ctx) {
      new Chart(ctx.getContext('2d'), {
        type: 'line',
        data: {
          labels: Array.from({ length: 24 }, (_, i) => `${i}:00`),
          datasets: [
            {
              label: 'CPU Usage',
              data: Array.from({ length: 24 }, () => Math.floor(Math.random() * 30) + 20),
              borderColor: '#4f46e5',
              backgroundColor: 'rgba(79, 70, 229, 0.1)',
              tension: 0.3,
              fill: true
            },
            {
              label: 'Memory Usage',
              data: Array.from({ length: 24 }, () => Math.floor(Math.random() * 30) + 40),
              borderColor: '#10b981',
              backgroundColor: 'rgba(16, 185, 129, 0.1)',
              tension: 0.3,
              fill: true
            }
          ]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'top',
            },
            tooltip: {
              mode: 'index',
              intersect: false,
            }
          },
          scales: {
            y: {
              beginAtZero: true,
              max: 100,
              ticks: {
                callback: (value) => `${value}%`
              }
            }
          }
        }
      });
    }
  }

  /**
   * Get color class based on usage percentage
   */
  getUsageColorClass(percent) {
    if (percent > 90) return 'bg-danger';
    if (percent > 70) return 'bg-warning';
    if (percent > 40) return 'bg-info';
    return 'bg-success';
  }

  /**
   * Format timestamp as time ago
   */
  formatTimeAgo(timestamp) {
    const seconds = Math.floor((new Date() - new Date(timestamp)) / 1000);
    
    const intervals = {
      year: 31536000,
      month: 2592000,
      week: 604800,
      day: 86400,
      hour: 3600,
      minute: 60,
      second: 1
    };
    
    for (const [unit, secondsInUnit] of Object.entries(intervals)) {
      const interval = Math.floor(seconds / secondsInUnit);
      if (interval >= 1) {
        return interval === 1 ? `1 ${unit} ago` : `${interval} ${unit}s ago`;
      }
    }
    
    return 'Just now';
  }

  /**
   * Get authentication headers
   */
  getAuthHeaders() {
    const token = this.getAuthToken();
    return {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : ''
    };
  }

  /**
   * Get authentication token from storage
   */
  getAuthToken() {
    // This should be implemented based on your auth system
    return localStorage.getItem('auth_token');
  }

  /**
   * Logout user
   */
  logout() {
    // Clear auth token
    localStorage.removeItem('auth_token');
    
    // Redirect to login page
    window.location.href = '/login';
  }

  /**
   * Show error message
   */
  showError(message) {
    // You can implement a toast or alert system here
    console.error('Error:', message);
    alert(message);
  }

  /**
   * Set loading state
   */
  setLoadingState(isLoading) {
    const buttons = document.querySelectorAll('button[type="button"]');
    buttons.forEach(button => {
      if (button.id !== 'logout') { // Don't disable logout button
        button.disabled = isLoading;
      }
    });
    
    // You can add a loading spinner or other UI indicators here
    if (isLoading) {
      document.body.style.cursor = 'wait';
    } else {
      document.body.style.cursor = 'default';
    }
  }
}

// Initialize the admin interface when the DOM is fully loaded
document.addEventListener('DOMContentLoaded', () => {
  window.mediaVaultAdmin = new MediaVaultAdmin();
});
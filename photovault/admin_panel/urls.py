from django.urls import path
from . import views
from .views import (
    DashboardView, UserListView, UserDetailView, UserCreateView, UserUpdateView,
    SystemSettingsView, SecurityLogsView, ContentModerationView, AuditLogsView,
    AnalyticsDashboardView, SystemHealthView, BackupRestoreView
)

app_name = 'admin_panel'

urlpatterns = [
    # Dashboard
    path('', DashboardView.as_view(), name='dashboard'),
    
    # User Management
    path('users/', UserListView.as_view(), name='user_list'),
    path('users/create/', UserCreateView.as_view(), name='user_create'),
    path('users/<int:pk>/', UserDetailView.as_view(), name='user_detail'),
    path('users/<int:pk>/edit/', UserUpdateView.as_view(), name='user_edit'),
    
    # System Settings
    path('settings/', SystemSettingsView.as_view(), name='system_settings'),
    
    # Security
    path('security/logs/', SecurityLogsView.as_view(), name='security_logs'),
    
    # Content Moderation
    path('moderation/', ContentModerationView.as_view(), name='content_moderation'),
    path('moderation/<int:pk>/review/', 
         views.ReviewContent.as_view(), 
         name='review_content'
    ),
    
    # Audit Logs
    path('audit-logs/', AuditLogsView.as_view(), name='audit_logs'),
    
    # Analytics
    path('analytics/', AnalyticsDashboardView.as_view(), name='analytics'),
    
    # System Health
    path('system/health/', SystemHealthView.as_view(), name='system_health'),
    
    # Backup & Restore
    path('backup/', BackupRestoreView.as_view(), name='backup_restore'),
    
    # API Endpoints
    path('api/users/', views.UserListAPIView.as_view(), name='api_users'),
    path('api/users/<int:pk>/', views.UserDetailAPIView.as_view(), name='api_user_detail'),
    path('api/security-events/', views.SecurityEventAPIView.as_view(), name='api_security_events'),
    path('api/audit-logs/', views.AuditLogAPIView.as_view(), name='api_audit_logs'),
    path('api/analytics/', views.AnalyticsAPIView.as_view(), name='api_analytics'),
]

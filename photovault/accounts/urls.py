from django.urls import path
from django.contrib.auth import views as auth_views
from . import views
from .views import (
    UserProfileView, UserProfileUpdateView, 
    UserPasswordChangeView, UserNotificationSettingsView,
    UserActivityLogView, UserDeleteView
)

app_name = 'accounts'

urlpatterns = [
    # User profile
    path('profile/', UserProfileView.as_view(), name='profile'),
    path('profile/edit/', UserProfileUpdateView.as_view(), name='profile_edit'),
    
    # Authentication
    path('login/', auth_views.LoginView.as_view(template_name='accounts/login.html'), name='login'),
    path('logout/', auth_views.LogoutView.as_view(next_page='home'), name='logout'),
    
    # Password management
    path('password/change/', UserPasswordChangeView.as_view(), name='password_change'),
    path('password/change/done/', 
         auth_views.PasswordChangeDoneView.as_view(
             template_name='accounts/password_change_done.html'
         ), 
         name='password_change_done'
    ),
    
    # Account settings
    path('settings/notifications/', 
         UserNotificationSettingsView.as_view(), 
         name='notification_settings'
    ),
    path('settings/security/', views.SecuritySettingsView.as_view(), name='security_settings'),
    
    # Activity and audit logs
    path('activity/', UserActivityLogView.as_view(), name='activity_log'),
    
    # Account deletion
    path('delete/', UserDeleteView.as_view(), name='delete_account'),
    
    # API endpoints
    path('api/profile/', views.ProfileAPIView.as_view(), name='api_profile'),
    path('api/activity/', views.ActivityLogAPIView.as_view(), name='api_activity'),
]

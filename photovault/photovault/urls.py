"""
URL configuration for photovault URL Configuration
"""
from django.conf import settings
from django.conf.urls.static import static
from django.contrib import admin
from django.urls import path, include
from django.views.generic import TemplateView
from django.contrib.auth import views as auth_views

# Admin site customization
admin.site.site_header = 'PhotoVault Administration'
admin.site.site_title = 'PhotoVault Admin'
admin.site.index_title = 'Welcome to PhotoVault Admin'

urlpatterns = [
    # Admin URLs
    path('admin/', admin.site.urls),
    
    # Authentication URLs (using Django Allauth)
    path('accounts/', include('allauth.urls')),
    
    # Custom account URLs
    path('accounts/', include('accounts.urls')),
    
    # Admin Panel URLs
    path('admin-panel/', include('admin_panel.urls', namespace='admin_panel')),
    
    # Partner Panel URLs
    path('partner/', include('partner_panel.urls', namespace='partner')),
    
    # Core URLs (API, etc.)
    path('api/', include('core.api.urls', namespace='api')),
    
    # Home page
    path('', TemplateView.as_view(template_name='home.html'), name='home'),
    
    # Password reset URLs (if not using Allauth's)
    path('password_reset/', auth_views.PasswordResetView.as_view(), name='password_reset'),
    path('password_reset/done/', auth_views.PasswordResetDoneView.as_view(), name='password_reset_done'),
    path('reset/<uidb64>/<token>/', auth_views.PasswordResetConfirmView.as_view(), name='password_reset_confirm'),
    path('reset/done/', auth_views.PasswordResetCompleteView.as_view(), name='password_reset_complete'),
]

# Serve media files in development
if settings.DEBUG:
    urlpatterns += static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
    urlpatterns += static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)

# Add debug toolbar if in development
if settings.DEBUG:
    import debug_toolbar
    urlpatterns = [
        path('__debug__/', include(debug_toolbar.urls)),
    ] + urlpatterns

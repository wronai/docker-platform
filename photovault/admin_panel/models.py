from django.db import models
from django.conf import settings
from django.utils.translation import gettext_lazy as _
from django.utils import timezone
from accounts.models import CustomUser


class SystemSetting(models.Model):
    """System-wide settings and configurations."""
    key = models.CharField(max_length=100, unique=True)
    value = models.TextField()
    description = models.TextField(blank=True, null=True)
    is_public = models.BooleanField(default=False)
    updated_at = models.DateTimeField(auto_now=True)
    updated_by = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True)
    
    def __str__(self):
        return f"{self.key}: {self.value[:50]}"


class SecurityEvent(models.Model):
    """Track security-related events and potential threats."""
    EVENT_TYPES = [
        ('login_attempt', 'Login Attempt'),
        ('failed_login', 'Failed Login'),
        ('password_reset', 'Password Reset'),
        ('suspicious_activity', 'Suspicious Activity'),
        ('brute_force', 'Brute Force Attempt'),
        ('account_lockout', 'Account Lockout'),
    ]
    
    event_type = models.CharField(max_length=50, choices=EVENT_TYPES)
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True)
    ip_address = models.GenericIPAddressField()
    user_agent = models.TextField(blank=True, null=True)
    details = models.JSONField(blank=True, null=True)
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.get_event_type_display()} - {self.ip_address} - {self.created_at}"


class ContentModeration(models.Model):
    """Track content that needs moderation."""
    STATUS_CHOICES = [
        ('pending', 'Pending Review'),
        ('approved', 'Approved'),
        ('rejected', 'Rejected'),
        ('flagged', 'Flagged for Review'),
    ]
    
    content_type = models.CharField(max_length=100)  # e.g., 'photo', 'comment', 'user_profile'
    content_id = models.PositiveIntegerField()
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='pending')
    reported_by = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True, related_name='reported_content')
    reviewed_by = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True, related_name='moderated_content')
    report_reason = models.TextField(blank=True, null=True)
    moderator_notes = models.TextField(blank=True, null=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    class Meta:
        ordering = ['-created_at']
        unique_together = ['content_type', 'content_id']
    
    def __str__(self):
        return f"{self.content_type} #{self.content_id} - {self.get_status_display()}"


class AuditLogEntry(models.Model):
    """Detailed audit log for admin actions."""
    ACTION_TYPES = [
        ('create', 'Create'),
        ('update', 'Update'),
        ('delete', 'Delete'),
        ('login', 'Login'),
        ('logout', 'Logout'),
        ('permission_change', 'Permission Change'),
        ('setting_change', 'Setting Change'),
    ]
    
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True)
    action = models.CharField(max_length=50, choices=ACTION_TYPES)
    model = models.CharField(max_length=100)
    object_id = models.CharField(max_length=100, blank=True, null=True)
    changes = models.JSONField(blank=True, null=True)
    ip_address = models.GenericIPAddressField(blank=True, null=True)
    user_agent = models.TextField(blank=True, null=True)
    timestamp = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['-timestamp']
        verbose_name_plural = 'Audit Log Entries'
    
    def __str__(self):
        return f"{self.get_action_display()} {self.model} by {self.user}"


class AdminNotification(models.Model):
    """Notifications for admin users."""
    NOTIFICATION_TYPES = [
        ('system_alert', 'System Alert'),
        ('user_report', 'User Report'),
        ('content_flagged', 'Content Flagged'),
        ('security_alert', 'Security Alert'),
        ('system_update', 'System Update'),
    ]
    
    notification_type = models.CharField(max_length=50, choices=NOTIFICATION_TYPES)
    title = models.CharField(max_length=200)
    message = models.TextField()
    is_read = models.BooleanField(default=False)
    related_content_type = models.CharField(max_length=100, blank=True, null=True)
    related_object_id = models.PositiveIntegerField(blank=True, null=True)
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.get_notification_type_display()}: {self.title}"

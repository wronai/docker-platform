from django.db import models
from django.conf import settings
from django.utils import timezone
from django.utils.translation import gettext_lazy as _


class TimeStampedModel(models.Model):
    """Abstract base class with self-updating created and modified fields."""
    created_at = models.DateTimeField(auto_now_add=True, db_index=True)
    updated_at = models.DateTimeField(auto_now=True, db_index=True)

    class Meta:
        abstract = True
        ordering = ['-created_at']


class ActivityLog(TimeStampedModel):
    """System-wide activity log"""
    ACTION_TYPES = [
        ('login', 'User Login'),
        ('logout', 'User Logout'),
        ('password_change', 'Password Changed'),
        ('profile_update', 'Profile Updated'),
        ('content_create', 'Content Created'),
        ('content_update', 'Content Updated'),
        ('content_delete', 'Content Deleted'),
        ('settings_update', 'Settings Updated'),
    ]
    
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True)
    action = models.CharField(max_length=50, choices=ACTION_TYPES)
    ip_address = models.GenericIPAddressField(blank=True, null=True)
    user_agent = models.TextField(blank=True, null=True)
    details = models.JSONField(blank=True, null=True)
    
    class Meta:
        verbose_name = _('Activity Log')
        verbose_name_plural = _('Activity Logs')
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.get_action_display()} - {self.user or 'System'}"


class SystemSetting(TimeStampedModel):
    """System-wide settings and configurations"""
    SETTING_TYPES = [
        ('string', 'String'),
        ('integer', 'Integer'),
        ('boolean', 'Boolean'),
        ('json', 'JSON'),
    ]
    
    key = models.CharField(max_length=100, unique=True, help_text="Setting key (snake_case)")
    value = models.TextField(blank=True, null=True, help_text="Setting value")
    setting_type = models.CharField(max_length=20, choices=SETTING_TYPES, default='string')
    is_public = models.BooleanField(default=False, help_text="If True, this setting can be exposed via API")
    description = models.TextField(blank=True, null=True, help_text="Description of the setting")
    updated_by = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True, 
                                 related_name='updated_settings')
    
    class Meta:
        verbose_name = _('System Setting')
        verbose_name_plural = _('System Settings')
        ordering = ['key']
    
    def __str__(self):
        return f"{self.key}: {self.value}"
    
    def get_typed_value(self):
        """Return the value cast to the appropriate type"""
        if not self.value:
            return None
            
        if self.setting_type == 'integer':
            try:
                return int(self.value)
            except (ValueError, TypeError):
                return 0
        elif self.setting_type == 'boolean':
            return self.value.lower() in ('true', '1', 't', 'y', 'yes')
        elif self.setting_type == 'json':
            import json
            try:
                return json.loads(self.value)
            except (json.JSONDecodeError, TypeError):
                return {}
        return self.value


class FileUpload(TimeStampedModel):
    """Generic file upload model"""
    UPLOAD_TYPES = [
        ('image', 'Image'),
        ('document', 'Document'),
        ('video', 'Video'),
        ('audio', 'Audio'),
        ('other', 'Other'),
    ]
    
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.SET_NULL, null=True, blank=True)
    file = models.FileField(upload_to='uploads/%Y/%m/%d/')
    file_type = models.CharField(max_length=20, choices=UPLOAD_TYPES)
    original_filename = models.CharField(max_length=255)
    mime_type = models.CharField(max_length=100, blank=True, null=True)
    file_size = models.PositiveIntegerField(help_text="File size in bytes")
    is_public = models.BooleanField(default=False)
    metadata = models.JSONField(blank=True, null=True)
    
    class Meta:
        verbose_name = _('File Upload')
        verbose_name_plural = _('File Uploads')
        ordering = ['-created_at']
    
    def __str__(self):
        return self.original_filename
    
    def get_file_size_display(self):
        """Return human-readable file size"""
        if self.file_size < 1024:
            return f"{self.file_size} B"
        elif self.file_size < 1024 * 1024:
            return f"{self.file_size / 1024:.1f} KB"
        elif self.file_size < 1024 * 1024 * 1024:
            return f"{self.file_size / (1024 * 1024):.1f} MB"
        return f"{self.file_size / (1024 * 1024 * 1024):.1f} GB"
    
    def get_absolute_url(self):
        return self.file.url if self.file else ''


class Notification(TimeStampedModel):
    """System notifications for users"""
    NOTIFICATION_TYPES = [
        ('info', 'Information'),
        ('success', 'Success'),
        ('warning', 'Warning'),
        ('error', 'Error'),
        ('system', 'System'),
    ]
    
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, related_name='notifications')
    notification_type = models.CharField(max_length=20, choices=NOTIFICATION_TYPES, default='info')
    title = models.CharField(max_length=200)
    message = models.TextField()
    is_read = models.BooleanField(default=False)
    read_at = models.DateTimeField(blank=True, null=True)
    action_url = models.URLField(blank=True, null=True)
    action_text = models.CharField(max_length=100, blank=True, null=True)
    
    class Meta:
        verbose_name = _('Notification')
        verbose_name_plural = _('Notifications')
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.get_notification_type_display()}: {self.title}"
    
    def mark_as_read(self, commit=True):
        self.is_read = True
        self.read_at = timezone.now()
        if commit:
            self.save()
    
    def mark_as_unread(self, commit=True):
        self.is_read = False
        self.read_at = None
        if commit:
            self.save()

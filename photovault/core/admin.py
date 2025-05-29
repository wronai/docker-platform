from django.contrib import admin
from django.utils.html import format_html
from django.urls import reverse
from django.utils.safestring import mark_safe
from django.conf import settings
from .models import (
    ActivityLog, SystemSetting, FileUpload, Notification
)


@admin.register(ActivityLog)
class ActivityLogAdmin(admin.ModelAdmin):
    list_display = ('action', 'user_display', 'ip_address', 'created_at')
    list_filter = ('action', 'created_at')
    search_fields = ('user__email', 'ip_address', 'action')
    readonly_fields = ('created_at', 'updated_at', 'details_formatted')
    date_hierarchy = 'created_at'
    
    def user_display(self, obj):
        if obj.user:
            url = reverse('admin:accounts_customuser_change', args=[obj.user.id])
            return mark_safe(f'<a href="{url}">{obj.user.email}</a>')
        return 'System'
    user_display.short_description = 'User'
    
    def details_formatted(self, obj):
        if not obj.details:
            return "No details available"
        return format_html('<pre>{}</pre>', str(obj.details))
    details_formatted.short_description = 'Details'
    
    def has_add_permission(self, request):
        return False
    
    def has_change_permission(self, request, obj=None):
        return False


@admin.register(SystemSetting)
class SystemSettingAdmin(admin.ModelAdmin):
    list_display = ('key', 'value_short', 'setting_type', 'is_public', 'updated_at')
    list_filter = ('setting_type', 'is_public')
    search_fields = ('key', 'value', 'description')
    list_editable = ('is_public',)
    readonly_fields = ('created_at', 'updated_at', 'updated_by_display', 'value_type_info')
    fieldsets = (
        ('Basic Information', {
            'fields': ('key', 'description', 'is_public')
        }),
        ('Value', {
            'fields': ('setting_type', 'value', 'value_type_info')
        }),
        ('Metadata', {
            'fields': ('updated_by_display', 'created_at', 'updated_at'),
            'classes': ('collapse',)
        }),
    )
    
    def value_short(self, obj):
        value = str(obj.value)
        return value[:100] + '...' if len(value) > 100 else value
    value_short.short_description = 'Value'
    
    def updated_by_display(self, obj):
        if obj.updated_by:
            url = reverse('admin:accounts_customuser_change', args=[obj.updated_by.id])
            return mark_safe(f'<a href="{url}">{obj.updated_by.email}</a>')
        return 'System'
    updated_by_display.short_description = 'Updated By'
    
    def value_type_info(self, obj):
        if obj.setting_type == 'boolean':
            return "Enter 'true' or 'false' (case insensitive)"
        elif obj.setting_type == 'integer':
            return "Enter a whole number"
        elif obj.setting_type == 'json':
            return "Enter valid JSON (e.g., {\"key\": \"value\"})"
        return "Text value"
    value_type_info.short_description = 'Format Hint'
    
    def save_model(self, request, obj, form, change):
        if not obj.pk:
            obj.updated_by = request.user
        super().save_model(request, obj, form, change)


@admin.register(FileUpload)
class FileUploadAdmin(admin.ModelAdmin):
    list_display = ('preview_thumbnail', 'original_filename', 'file_type', 'user_display', 'file_size_display', 'created_at')
    list_filter = ('file_type', 'created_at')
    search_fields = ('original_filename', 'user__email')
    readonly_fields = ('created_at', 'updated_at', 'preview_large', 'file_size_display', 'metadata_formatted')
    list_select_related = ('user',)
    
    fieldsets = (
        ('File Information', {
            'fields': ('preview_large', 'original_filename', 'file', 'file_type', 'is_public')
        }),
        ('Metadata', {
            'fields': ('mime_type', 'file_size_display', 'metadata_formatted'),
            'classes': ('collapse',)
        }),
        ('Ownership', {
            'fields': ('user',)
        }),
        ('Timestamps', {
            'fields': ('created_at', 'updated_at'),
            'classes': ('collapse',)
        }),
    )
    
    def preview_thumbnail(self, obj):
        if obj.file and obj.file_type == 'image':
            return mark_safe(f'<img src="{obj.file.url}" style="max-height: 50px; max-width: 50px;" />')
        return ""
    preview_thumbnail.short_description = 'Preview'
    
    def preview_large(self, obj):
        if obj.file and obj.file_type == 'image':
            return mark_safe(f'<img src="{obj.file.url}" style="max-height: 300px; max-width: 100%;" />')
        return "No preview available"
    preview_large.short_description = 'Preview'
    
    def file_size_display(self, obj):
        return obj.get_file_size_display()
    file_size_display.short_description = 'File Size'
    
    def user_display(self, obj):
        if obj.user:
            url = reverse('admin:accounts_customuser_change', args=[obj.user.id])
            return mark_safe(f'<a href="{url}">{obj.user.email}</a>')
        return 'System'
    user_display.short_description = 'User'
    
    def metadata_formatted(self, obj):
        if not obj.metadata:
            return "No metadata available"
        return format_html('<pre>{}</pre>', str(obj.metadata))
    metadata_formatted.short_description = 'Metadata'


@admin.register(Notification)
class NotificationAdmin(admin.ModelAdmin):
    list_display = ('title', 'user_display', 'notification_type', 'is_read', 'created_at')
    list_filter = ('notification_type', 'is_read', 'created_at')
    search_fields = ('title', 'message', 'user__email')
    list_editable = ('is_read',)
    readonly_fields = ('created_at', 'updated_at', 'read_at_display')
    actions = ['mark_as_read', 'mark_as_unread']
    
    fieldsets = (
        ('Notification', {
            'fields': ('user', 'title', 'message', 'notification_type')
        }),
        ('Actions', {
            'fields': ('is_read', 'read_at_display', 'action_url', 'action_text')
        }),
        ('Timestamps', {
            'fields': ('created_at', 'updated_at'),
            'classes': ('collapse',)
        }),
    )
    
    def user_display(self, obj):
        url = reverse('admin:accounts_customuser_change', args=[obj.user.id])
        return mark_safe(f'<a href="{url}">{obj.user.email}</a>')
    user_display.short_description = 'User'
    
    def read_at_display(self, obj):
        if obj.read_at:
            return obj.read_at.strftime('%Y-%m-%d %H:%M:%S')
        return 'Not read'
    read_at_display.short_description = 'Read At'
    
    def mark_as_read(self, request, queryset):
        updated = queryset.update(is_read=True, read_at=timezone.now())
        self.message_user(request, f"Marked {updated} notifications as read.")
    mark_as_read.short_description = "Mark selected notifications as read"
    
    def mark_as_unread(self, request, queryset):
        updated = queryset.update(is_read=False, read_at=None)
        self.message_user(request, f"Marked {updated} notifications as unread.")
    mark_as_unread.short_description = "Mark selected notifications as unread"
    
    def get_queryset(self, request):
        return super().get_queryset(request).select_related('user')

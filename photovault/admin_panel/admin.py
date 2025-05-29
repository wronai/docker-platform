from django.contrib import admin
from django.utils.html import format_html
from .models import (
    SystemSetting, SecurityEvent, ContentModeration, 
    AuditLogEntry, AdminNotification, SystemConfig, AuditLog
)


@admin.register(SystemSetting)
class SystemSettingAdmin(admin.ModelAdmin):
    list_display = ('key', 'value_short', 'is_public', 'updated_at', 'updated_by')
    list_filter = ('is_public', 'updated_at')
    search_fields = ('key', 'value', 'description')
    list_editable = ('is_public',)
    readonly_fields = ('updated_at', 'updated_by')
    
    def value_short(self, obj):
        return obj.value[:100] + '...' if len(obj.value) > 100 else obj.value
    value_short.short_description = 'Value'
    
    def save_model(self, request, obj, form, change):
        if not obj.pk:
            obj.updated_by = request.user
        super().save_model(request, obj, form, change)


@admin.register(SecurityEvent)
class SecurityEventAdmin(admin.ModelAdmin):
    list_display = ('event_type', 'user_display', 'ip_address', 'created_at')
    list_filter = ('event_type', 'created_at')
    search_fields = ('user__email', 'ip_address', 'details')
    readonly_fields = ('created_at', 'details_formatted')
    date_hierarchy = 'created_at'
    
    def user_display(self, obj):
        return obj.user.email if obj.user else 'System'
    user_display.short_description = 'User'
    
    def details_formatted(self, obj):
        if not obj.details:
            return "No details available"
        return format_html('<pre>{}</pre>', str(obj.details))
    details_formatted.short_description = 'Event Details'


class ContentModerationInline(admin.TabularInline):
    model = ContentModeration
    extra = 0
    readonly_fields = ('created_at', 'updated_at', 'reported_by', 'reviewed_by')
    can_delete = False


@admin.register(ContentModeration)
class ContentModerationAdmin(admin.ModelAdmin):
    list_display = ('content_type', 'object_id', 'status', 'moderated_by', 'moderated_at')
    list_filter = ('status', 'content_type')
    search_fields = ('object_id', 'moderated_by__email', 'notes')
    readonly_fields = ('moderated_at',)
    actions = ['approve_content', 'reject_content']

    def approve_content(self, request, queryset):
        queryset.update(status='approved', moderated_by=request.user)
    approve_content.short_description = "Approve selected content"

    def reject_content(self, request, queryset):
        queryset.update(status='rejected', moderated_by=request.user)
    reject_content.short_description = "Reject selected content"
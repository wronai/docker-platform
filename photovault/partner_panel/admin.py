from django.contrib import admin
from django.utils.html import format_html
from django.urls import reverse
from django.utils.safestring import mark_safe
from .models import (
    PhotoBatch, Photo, PhotoShare, PhotoAccessLog, 
    PhotoCollection, PhotoComment
)


class PhotoInline(admin.TabularInline):
    model = Photo
    extra = 0
    readonly_fields = ('preview_thumbnail', 'file_size_display', 'status')
    fields = ('preview_thumbnail', 'title', 'status', 'is_featured', 'created_at')
    
    def preview_thumbnail(self, obj):
        if obj.image:
            return mark_safe(f'<img src="{obj.image.url}" style="max-height: 100px; max-width: 100px;" />')
        return "No Image"
    preview_thumbnail.short_description = 'Preview'
    
    def file_size_display(self, obj):
        if obj.file_size < 1024:
            return f"{obj.file_size} B"
        elif obj.file_size < 1024 * 1024:
            return f"{obj.file_size / 1024:.1f} KB"
        return f"{obj.file_size / (1024 * 1024):.1f} MB"
    file_size_display.short_description = 'File Size'
    
    def has_change_permission(self, request, obj=None):
        return False


@admin.register(PhotoBatch)
class PhotoBatchAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'created_at', 'status')
    list_filter = ('status', 'created_at')
    search_fields = ('user__email', 'id')
    inlines = [PhotoInline]
    readonly_fields = ('created_at', 'updated_at')

@admin.register(Photo)
class PhotoAdmin(admin.ModelAdmin):
    list_display = ('id', 'batch', 'user', 'created_at', 'is_approved')
    list_filter = ('is_approved', 'created_at')
    search_fields = ('description', 'user__email', 'batch__id')
    readonly_fields = ('created_at', 'updated_at')
    list_editable = ('is_approved',)

@admin.register(PhotoShare)
class PhotoShareAdmin(admin.ModelAdmin):
    list_display = ('id', 'photo', 'shared_by', 'shared_with', 'can_edit', 'created_at')
    list_filter = ('can_edit', 'created_at')
    search_fields = ('photo__id', 'shared_by__email', 'shared_with__email')

@admin.register(PhotoAccess)
class PhotoAccessAdmin(admin.ModelAdmin):
    list_display = ('id', 'user', 'photo', 'access_type', 'granted_at')
    list_filter = ('access_type', 'granted_at')
    search_fields = ('user__email', 'photo__id')
from django.contrib import admin
from .models import Photo, PhotoBatch, PhotoShare, PhotoAccess

class PhotoInline(admin.TabularInline):
    model = Photo
    extra = 1
    fields = ('image', 'description', 'is_approved')
    readonly_fields = ('created_at',)

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
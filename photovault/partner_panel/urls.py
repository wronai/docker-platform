from django.urls import path
from . import views
from .views import (
    PartnerDashboardView, PhotoBatchListView, PhotoBatchCreateView,
    PhotoBatchDetailView, PhotoBatchUpdateView, PhotoUploadView,
    PhotoEditView, BatchShareView, BatchAnalyticsView, BatchExportView,
    BatchSettingsView, BatchMetadataView, BatchActionsView, BatchShareSettingsView
)

app_name = 'partner_panel'

urlpatterns = [
    # Dashboard
    path('', PartnerDashboardView.as_view(), name='dashboard'),
    
    # Photo Batches
    path('batches/', PhotoBatchListView.as_view(), name='batch_list'),
    path('batches/create/', PhotoBatchCreateView.as_view(), name='batch_create'),
    path('batches/<int:pk>/', PhotoBatchDetailView.as_view(), name='batch_detail'),
    path('batches/<int:pk>/edit/', PhotoBatchUpdateView.as_view(), name='batch_edit'),
    path('batches/<int:pk>/upload/', PhotoUploadView.as_view(), name='batch_upload'),
    
    # Batch Actions
    path('batches/<int:pk>/share/', BatchShareView.as_view(), name='batch_share'),
    path('batches/<int:pk>/analytics/', BatchAnalyticsView.as_view(), name='batch_analytics'),
    path('batches/<int:pk>/export/', BatchExportView.as_view(), name='batch_export'),
    path('batches/<int:pk>/settings/', BatchSettingsView.as_view(), name='batch_settings'),
    path('batches/<int:pk>/metadata/', BatchMetadataView.as_view(), name='batch_metadata'),
    path('batches/<int:pk>/actions/', BatchActionsView.as_view(), name='batch_actions'),
    path('batches/<int:pk>/share-settings/', 
         BatchShareSettingsView.as_view(), 
         name='batch_share_settings'
    ),
    
    # Photos
    path('photos/<int:pk>/edit/', PhotoEditView.as_view(), name='photo_edit'),
    
    # API Endpoints
    path('api/batches/', views.BatchListAPIView.as_view(), name='api_batch_list'),
    path('api/batches/<int:pk>/', views.BatchDetailAPIView.as_view(), name='api_batch_detail'),
    path('api/batches/<int:batch_id>/photos/', views.PhotoListAPIView.as_view(), name='api_photo_list'),
    path('api/batches/<int:batch_id>/photos/<int:pk>/', 
         views.PhotoDetailAPIView.as_view(), 
         name='api_photo_detail'
    ),
    path('api/shared-links/', views.SharedLinkAPIView.as_view(), name='api_shared_links'),
    path('api/analytics/', views.PartnerAnalyticsAPIView.as_view(), name='api_partner_analytics'),
]

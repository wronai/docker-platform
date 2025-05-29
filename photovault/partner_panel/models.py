from django.db import models
from django.conf import settings
from django.utils import timezone

class PhotoBatch(models.Model):
    STATUS_CHOICES = [
        ('uploading', 'Uploading'),
        ('uploaded', 'Uploaded'),
        ('processing', 'Processing'),
        ('processed', 'Processed'),
        ('failed', 'Failed'),
    ]
    
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE)
    status = models.CharField(max_length=20, choices=STATUS_CHOICES, default='uploading')
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
        
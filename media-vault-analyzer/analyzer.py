import os
import time
import requests
import cv2
import numpy as np
from PIL import Image, ExifTags
import tensorflow as tf
from transformers import BlipProcessor, BlipForConditionalGeneration
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class MediaVaultAnalyzer:
    def __init__(self):
        self.api_url = os.getenv('API_URL', 'http://media-vault-api:8080')
        self.processing_path = os.getenv('PROCESSING_PATH', '/processing')
        self.output_path = os.getenv('OUTPUT_PATH', '/uploads')

        # Initialize AI models
        self.load_models()

    def load_models(self):
        """Load AI models for image analysis"""
        try:
            # Load image captioning model
            logger.info("Loading BLIP model for image captioning...")
            self.processor = BlipProcessor.from_pretrained("Salesforce/blip-image-captioning-base")
            self.caption_model = BlipForConditionalGeneration.from_pretrained("Salesforce/blip-image-captioning-base")
            logger.info("BLIP model loaded successfully")

        except Exception as e:
            logger.error(f"Failed to load models: {e}")
            self.processor = None
            self.caption_model = None

    def generate_description(self, image_path):
        """Generate AI description for an image"""
        if not self.processor or not self.caption_model:
            return "AI description service unavailable"

        try:
            # Load and preprocess image
            image = Image.open(image_path).convert('RGB')

            # Generate caption
            inputs = self.processor(image, return_tensors="pt")
            out = self.caption_model.generate(**inputs, max_length=50)
            caption = self.processor.decode(out[0], skip_special_tokens=True)

            return caption

        except Exception as e:
            logger.error(f"Failed to generate description: {e}")
            return "Failed to generate description"

    def extract_metadata(self, image_path):
        """Extract metadata from image"""
        metadata = {}

        try:
            # Open image with PIL
            image = Image.open(image_path)

            # Get basic info
            metadata['width'] = image.width
            metadata['height'] = image.height
            metadata['format'] = image.format
            metadata['mode'] = image.mode

            # Extract EXIF data
            exif_data = {}
            if hasattr(image, '_getexif') and image._getexif():
                exif = image._getexif()
                for tag_id, value in exif.items():
                    tag = ExifTags.TAGS.get(tag_id, tag_id)
                    exif_data[tag] = str(value)

            metadata['exif'] = exif_data

            # Camera info
            if 'Make' in exif_data:
                metadata['camera_make'] = exif_data['Make']
            if 'Model' in exif_data:
                metadata['camera_model'] = exif_data['Model']
            if 'DateTime' in exif_data:
                metadata['taken_at'] = exif_data['DateTime']

        except Exception as e:
            logger.error(f"Failed to extract metadata: {e}")

        return metadata

    def generate_thumbnail(self, image_path, photo_id):
        """Generate thumbnail for image"""
        try:
            thumbnail_dir = os.path.join(self.output_path, 'thumbnails')
            os.makedirs(thumbnail_dir, exist_ok=True)

            # Load image
            image = Image.open(image_path)

            # Create thumbnail
            thumbnail_size = (300, 300)
            image.thumbnail(thumbnail_size, Image.Resampling.LANCZOS)

            # Save thumbnail
            thumbnail_path = os.path.join(thumbnail_dir, f"{photo_id}_thumb.jpg")
            image.save(thumbnail_path, "JPEG", quality=85)

            return thumbnail_path

        except Exception as e:
            logger.error(f"Failed to generate thumbnail: {e}")
            return None

    def process_photo(self, photo_data):
        """Process a single photo"""
        photo_id = photo_data['id']
        file_path = photo_data['file_path']

        logger.info(f"Processing photo {photo_id}: {file_path}")

        try:
            # Generate AI description
            ai_description = self.generate_description(file_path)

            # Extract metadata
            metadata = self.extract_metadata(file_path)

            # Generate thumbnail
            thumbnail_path = self.generate_thumbnail(file_path, photo_id)

            # Update photo via API
            update_data = {
                'ai_description': ai_description,
                'ai_confidence': 0.85,  # Mock confidence score
                'metadata': metadata,
                'thumbnail_path': thumbnail_path,
                'processed_at': time.strftime('%Y-%m-%d %H:%M:%S'),
                'moderation_status': 'approved'  # Basic approval
            }

            # Send update to API
            response = requests.put(
                f"{self.api_url}/api/vault/files/{photo_id}",
                json=update_data,
                timeout=30
            )

            if response.status_code == 200:
                logger.info(f"Successfully processed photo {photo_id}")
            else:
                logger.error(f"Failed to update photo {photo_id}: {response.status_code}")

        except Exception as e:
            logger.error(f"Failed to process photo {photo_id}: {e}")

    def run(self):
        """Main processing loop"""
        logger.info("Media Vault Analyzer started")

        while True:
            try:
                # Get pending photos from API
                response = requests.get(
                    f"{self.api_url}/api/admin/photos/pending",
                    timeout=30
                )

                if response.status_code == 200:
                    pending_photos = response.json()

                    for photo in pending_photos:
                        self.process_photo(photo)
                        time.sleep(1)  # Small delay between processing

                # Wait before next check
                time.sleep(10)

            except Exception as e:
                logger.error(f"Error in main loop: {e}")
                time.sleep(30)  # Wait longer on error

if __name__ == "__main__":
    analyzer = MediaVaultAnalyzer()
    analyzer.run()


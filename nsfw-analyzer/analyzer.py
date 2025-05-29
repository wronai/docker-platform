import os
import requests
import numpy as np
import tensorflow as tf
from PIL import Image
from flask import Flask, request, jsonify
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

class NSFWAnalyzer:
    def __init__(self):
        self.model_path = os.getenv('MODEL_PATH', '/models')
        self.confidence_threshold = float(os.getenv('CONFIDENCE_THRESHOLD', '0.8'))
        self.model = None
        self.load_model()

    def load_model(self):
        """Load NSFW detection model"""
        try:
            # Try to load pre-trained model
            model_file = os.path.join(self.model_path, 'nsfw_mobilenet.h5')
            if os.path.exists(model_file):
                self.model = tf.keras.models.load_model(model_file)
                logger.info("NSFW model loaded successfully")
            else:
                logger.warning("NSFW model not found, using mock detection")
                self.model = None
        except Exception as e:
            logger.error(f"Failed to load NSFW model: {e}")
            self.model = None

    def preprocess_image(self, image_path):
        """Preprocess image for NSFW detection"""
        try:
            image = Image.open(image_path).convert('RGB')
            image = image.resize((224, 224))
            image_array = np.array(image) / 255.0
            image_array = np.expand_dims(image_array, axis=0)
            return image_array
        except Exception as e:
            logger.error(f"Failed to preprocess image: {e}")
            return None

    def detect_nsfw(self, image_path):
        """Detect NSFW content in image"""
        if not self.model:
            # Mock detection for development
            return {
                'is_nsfw': False,
                'confidence': 0.1,
                'categories': {
                    'safe': 0.9,
                    'suggestive': 0.05,
                    'explicit': 0.05
                }
            }

        try:
            # Preprocess image
            image_array = self.preprocess_image(image_path)
            if image_array is None:
                return None

            # Run prediction
            predictions = self.model.predict(image_array)

            # Parse results (assuming model outputs probabilities for each category)
            categories = {
                'safe': float(predictions[0][0]),
                'suggestive': float(predictions[0][1]),
                'explicit': float(predictions[0][2])
            }

            # Determine if NSFW
            nsfw_confidence = categories['suggestive'] + categories['explicit']
            is_nsfw = nsfw_confidence > self.confidence_threshold

            return {
                'is_nsfw': is_nsfw,
                'confidence': nsfw_confidence,
                'categories': categories
            }

        except Exception as e:
            logger.error(f"NSFW detection failed: {e}")
            return None

# Global analyzer instance
analyzer = NSFWAnalyzer()

@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({'status': 'healthy', 'model_loaded': analyzer.model is not None})

@app.route('/analyze', methods=['POST'])
def analyze_image():
    """Analyze image for NSFW content"""
    try:
        data = request.get_json()
        image_path = data.get('image_path')

        if not image_path or not os.path.exists(image_path):
            return jsonify({'error': 'Invalid image path'}), 400

        # Run NSFW detection
        result = analyzer.detect_nsfw(image_path)

        if result is None:
            return jsonify({'error': 'Analysis failed'}), 500

        return jsonify(result)

    except Exception as e:
        logger.error(f"Analysis request failed: {e}")
        return jsonify({'error': 'Internal server error'}), 500

@app.route('/batch-analyze', methods=['POST'])
def batch_analyze():
    """Analyze multiple images"""
    try:
        data = request.get_json()
        image_paths = data.get('image_paths', [])

        results = []
        for image_path in image_paths:
            if os.path.exists(image_path):
                result = analyzer.detect_nsfw(image_path)
                results.append({
                    'image_path': image_path,
                    'result': result
                })
            else:
                results.append({
                    'image_path': image_path,
                    'error': 'File not found'
                })

        return jsonify({'results': results})

    except Exception as e:
        logger.error(f"Batch analysis failed: {e}")
        return jsonify({'error': 'Internal server error'}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8501, debug=False)



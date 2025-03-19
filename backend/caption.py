import sys
import json
from transformers import BlipProcessor, BlipForConditionalGeneration
from PIL import Image

# Load the BLIP captioning model
processor = BlipProcessor.from_pretrained("Salesforce/blip-image-captioning-base")
model = BlipForConditionalGeneration.from_pretrained("Salesforce/blip-image-captioning-base")

# Read image path from arguments
image_path = sys.argv[1]

try:
    image = Image.open(image_path).convert("RGB")

    # Process and generate caption
    inputs = processor(image, return_tensors="pt")
    caption_ids = model.generate(**inputs)
    caption = processor.batch_decode(caption_ids, skip_special_tokens=True)[0]

    # Print response as JSON
    print(json.dumps({"caption": caption}))

except Exception as e:
    print(json.dumps({"error": str(e)}))

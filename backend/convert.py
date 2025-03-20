import sys
import json
from transformers import pipeline

def generate_poetry(caption):
    generator = pipeline("text-generation", model="gpt2")
    prompt = f"Turn this caption into a poetic verse:\n{caption}\nPoem:"
    result = generator(prompt, max_length=100, num_return_sequences=1)
    return result[0]['generated_text'].split("Poem:")[-1].strip()

if __name__ == "__main__":
    # Read input from Go backend
    input_json = sys.stdin.read()
    data = json.loads(input_json)
    
    caption = data.get("caption", "")
    if not caption:
        print(json.dumps({"error": "No caption provided"}))
        sys.exit(1)
    
    poem = generate_poetry(caption)
    
    # Return JSON response
    print(json.dumps({"poem": poem}))

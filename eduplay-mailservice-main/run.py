from mail import app
from dotenv import load_dotenv
import os

load_dotenv()
PORT = os.getenv('PORT')

if __name__ == "__main__":
    app.run(hort="0.0.0.0", port=int(PORT), debug=True)

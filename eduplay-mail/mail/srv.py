from flask import Flask, request
import os
from dotenv import load_dotenv
from .email import Emailer


load_dotenv()
PORT = os.getenv('PORT')
EM_USER = os.getenv('EM_USER')
EM_PASSWD = os.getenv('EM_PASSWORD')
EM_SENDERNAME = os.getenv("EM_SENDERNAME")
EM_SENDEREMAIL = os.getenv("EM_SENDEREMAIL")
EM_SRV = os.getenv('EM_SRV')
EM_PORT = os.getenv('EM_PORT')


app = Flask(__name__)
emailer = Emailer(EM_USER, EM_PASSWD, 
                  EM_SRV, EM_PORT, 
                  sender_name=EM_SENDERNAME, 
                  sender_email=EM_SENDEREMAIL)

# @app.route("/send", methods=["POST"])
# def sendemail():
#     """
#     Docstring для sendemail

#     Expected format in POST request:
#     {
#         "to": rec | [rec1, rec2],   # reciever(s)
#         "subject": <subject>,   # email subject
#         "body": <HTML-formatted body>,  # html text body
#         ["cc": rec | [rec1, rec2] # copy reciever(s) ],
#         ["broadcast": True | False (True by default)]
#     }

#     """

#     if request.form is None:
#         return "Bad request", 400
    
#     req_keys = set(["to", "subject", "body"])
#     if len(req_keys.intersection(set(request.form.keys()))) != len(req_keys):
#         return "Missing fields", 400
    
#     recievers = [request.form["to"]] if type(request.form["to"]) is str else request.form["to"]
#     sbj = request.form["subject"]
#     body = request.form["body"]
#     cc_recievers = request.form.get("cc", [])
#     cc_recievers = [cc_recievers] if type(cc_recievers) is str else cc_recievers
#     do_broadcast = request.form.get("broadcast", True)

#     if do_broadcast:
#         emailer.send(recievers, sbj, body, 
#                     copy_receivers=cc_recievers, 
#                     broadcast=do_broadcast)
#     else:
#         for rec in recievers:
#             emailer.send([rec], sbj, body, 
#                     copy_receivers=cc_recievers, 
#                     broadcast=do_broadcast)

#     return "ok"


@app.route("/send", methods=["POST"])
def sendemail():
    """
    Docstring для sendemail

    Expected format in POST request:
    {
        "to": rec | [rec1, rec2],   # reciever(s)
        "subject": <subject>,   # email subject
        "body": <HTML-formatted body>,  # html text body
        ["cc": rec | [rec1, rec2] # copy reciever(s) ],
        ["broadcast": True | False (True by default)]
    }

    """
    data = request.get_json()

    if not data:
        return "Bad request", 400

    req_keys = {"to", "subject", "body"}
    if not req_keys.issubset(data.keys()):
        return "Missing fields", 400

    receivers = data["to"]
    if isinstance(receivers, str):
        receivers = [receivers]

    subject = data["subject"]
    body = data["body"]

    cc = data.get("cc", [])
    if isinstance(cc, str):
        cc = [cc]

    broadcast = data.get("broadcast", True)

    if broadcast:
        emailer.send(receivers, subject, body, copy_receivers=cc)
    else:
        for rec in receivers:
            emailer.send([rec], subject, body, copy_receivers=cc)

    return "ok", 200
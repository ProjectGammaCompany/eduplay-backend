import smtplib
import ssl
import certifi

from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.application import MIMEApplication
from email.header import Header
from os.path import basename
from email.utils import formataddr


class Emailer:
    def __init__(self, email, passwd, server_name, server_port, sender_name=None, sender_email=None):
        self._send_server = smtplib.SMTP(server_name, server_port)
        self.__email = email
        self.__passwd = passwd
        self._srv = server_name
        self._port = server_port

        self._login()

        if sender_name is None:
            sender_name = email
        if sender_email is None:
            sender_email = email
        self.sender = [sender_name, sender_email]

    def send(self, receivers, subject, body, attachments=None, copy_receivers=None, broadcast=True):
        message = MIMEMultipart()
        message["From"] = formataddr((str(Header(self.sender[0], 'utf-8')),
                                      self.sender[1]))
        message["To"] = ",".join(receivers)
        if copy_receivers is not None:
            message["Cc"] = ",".join(copy_receivers)
        message["Subject"] = subject
        message.attach(MIMEText(body, "html"))

        if attachments is not None:
            for fpath in attachments:
                with open(fpath, "rb") as fil:
                    part = MIMEApplication(
                        fil.read(),
                        Name=basename(fpath)
                    )
                part['Content-Disposition'] = 'attachment; filename="%s"' % basename(fpath)
                message.attach(part)
        if copy_receivers is None:
            copy_receivers = []
        
        if not self._is_connected():
            self._login()
        
        self._send_server.sendmail(self.sender[1],
                                   receivers + copy_receivers,
                                   message.as_string())
        
        # if broadcast:
        #     self._send_server.sendmail(self.sender[1],
        #                                 receivers + copy_receivers,
        #                                 message.as_string())
        # else:
        #     for rec in receivers:
        #         message["To"] = rec
        #         self._send_server.sendmail(self.sender[1],
        #                                 [rec] + copy_receivers,
        #                                 message.as_string())
        # self._logout()

        return 200, "OK"
    
    def _login(self):
        self._send_server.connect(self._srv, self._port)
        context = ssl.create_default_context()
        context.load_verify_locations(certifi.where())
        self._send_server.ehlo()  # Identify yourself to the server
        self._send_server.starttls(context=context) # Secure the connection
        self._send_server.login(self.__email, self.__passwd)

    def _logout(self):
        self._send_server.quit()

    def _is_connected(self):
        try:
            status = self._send_server.noop()[0]
            return status == 250
        except smtplib.SMTPServerDisconnected:
            return False
        except Exception:
            return False

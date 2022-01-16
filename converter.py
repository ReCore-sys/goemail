import os
import uuid

for x in os.listdir(r"C:/Users/ReCor/Documents/OtherCode/gomail/static/img"):
    uid = uuid.uuid4().hex
    cmd = r"rename C:\Users\ReCor\Documents\OtherCode\gomail\static\img" + \
        "\\" + x + " " + uid.replace('-', '') + ".png"
    os.system(cmd)

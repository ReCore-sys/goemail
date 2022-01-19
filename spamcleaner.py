

with open(r"C:\Users\ReCor\Documents\OtherCode\gomail\spams.txt") as f:
    lines = f.readlines()
    lines = [x.strip() for x in lines]
    lines.sort()
    print(lines)
with open(r"C:\Users\ReCor\Documents\OtherCode\gomail\spams.txt", "w") as f:
    f.write("".join(lines))

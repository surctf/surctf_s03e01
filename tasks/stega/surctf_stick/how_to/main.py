import os

in_ffconcat = """ffconcat version 1.0
file ./images/%s
duration 240
file true_flag.png"""

ffmpeg_cmd = "ffmpeg -safe 0 -hwaccel cuda -i in.ffconcat -r 1 -vcodec libvpx-vp9 -pix_fmt yuva420p -s 512x512 -b:v 0 -crf 30 -pass 2 -cpu-used -1  ./result/%s.webm"
tgradish_cdm = "tgradish spoof %s %s"

files = os.listdir("./images")
files.sort()

for fname in files:
    with open("in.ffconcat", "w") as f:
        f.write(in_ffconcat % fname)

    name, ext = fname.split(".")
    os.system(ffmpeg_cmd % name)
    os.system(tgradish_cdm % ("./result/%s.webm" % name, "./result/spoof_%s.webm" % name))
    print(fname)
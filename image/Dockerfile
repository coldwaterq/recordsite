FROM selenium/standalone-firefox-debug

RUN sudo apt-get update
RUN sudo apt-get install pulseaudio -y
RUN sudo apt-get install software-properties-common -y
RUN sudo add-apt-repository ppa:jonathonf/ffmpeg-4
RUN sudo apt-get install ffmpeg -y
RUN sudo apt-get install alsa -y
RUN mkdir /tmp/recordings
RUN sudo apt-get install vlc -y
COPY watch.mp4 /tmp/
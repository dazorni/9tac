FROM debian:wheezy

EXPOSE 5000

COPY build/ /
ADD build/9tac.x64.1.5 /9tac
CMD /9tac

FROM public.ecr.aws/docker/library/alpine:3.20

WORKDIR /
ADD gogame /gogame

CMD ["/gogame"]

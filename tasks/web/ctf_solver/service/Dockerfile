FROM python:3.9-slim-buster

WORKDIR app

RUN apt update
RUN apt install -y libpq-dev gcc

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY *.py ./
COPY .env .
ADD help ./help/

CMD ["gunicorn","--bind", "0.0.0.0", "-w", "4", "main:app"]
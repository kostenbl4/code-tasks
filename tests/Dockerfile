FROM python:3.9-slim

WORKDIR /app

COPY tests/* .

RUN pip install -r requirements.txt

CMD ["pytest", "tests_all.py", "-v"]
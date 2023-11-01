FROM python:3.11

# Set work directory
WORKDIR /app
RUN mkdir /app/logs

# Copy project
COPY . .

# Install dependencies
RUN pip install --upgrade pip
RUN pip3 install -r requirements.txt

# Command to run on container start
ENTRYPOINT [ "python", "-m", "streamservice.main" ]
CMD [ "bitstamp" ]

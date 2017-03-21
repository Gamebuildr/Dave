FROM debian:squeeze

# Add Compiled Dave 
RUN mkdir -p /var/www/go/bin; \
    mkdir -p /var/www/go/big/client/logs
COPY Dave /var/www/go/bin

# Run Dave
CMD /var/www/go/bin/Dave

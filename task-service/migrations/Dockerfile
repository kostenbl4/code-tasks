FROM gomicro/goose

ADD ./task-service/migrations/*.sql /migrations/
ADD ./task-service/migrations/entrypoint.sh /migrations/

ENTRYPOINT ["sh", "/migrations/entrypoint.sh"]
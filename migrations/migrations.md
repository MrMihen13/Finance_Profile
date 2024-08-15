# Migrations

## Описание

Миграции создаются при помощи sql скриптов, предназначенных для запуска в базе данных. Первоначальный скрипт
содержит логику документирования миграций и имеет название `0_<YYYYMMDDHHMM>_init.sql`.

Скрипт выполняет слудующую логику:

1. Создает таблице хранения версий миграций;

```postgresql
   CREATE TABLE IF NOT EXISTS "migrations"
   (
       "id"          integer PRIMARY KEY,
       "migrated_at" timestamptz NOT NULL DEFAULT now()
   );
```

2. Добовляет функицю логирования версии миграции;

```postgresql
   CREATE OR REPLACE FUNCTION log_migration(id integer) RETURNS void AS
$$
BEGIN
    INSERT INTO "migrations" (id) VALUES (id);
END;
$$ LANGUAGE plpgsql;
```

3. Добовляет функцию проверки последней версии миграции.

```postgresql
   CREATE OR REPLACE FUNCTION assert_latest_migration(id integer) RETURNS void AS
$$
DECLARE
    latest_id integer;
BEGIN
    SELECT MAX(migrations.id) INTO latest_id FROM migrations;

    ASSERT latest_id = id, 'migration assertion ' || id || ' failed, current latest is ' || latest_id;
    RETURN;
END;
$$ LANGUAGE plpgsql;
```

## Правила создания

### Название миграции

Мигрции должны иметь следующий вид `<номер>_<дата>_<название>.sql`, где:
* **Номер** - номер миграции отличающийся на 1 от предыдущей миграции,
* **Дата** - дата и время в формате `<YYYYMMDD><HHMM>`
* **Название** - краткое описание миграции.

_Например_: Миграция под номером 16 созданная 25 декабря 2012 года в 12:54 для добавления таблицы User, будет иметь 
следующий вид:

`16_201212251254_add_user_table.sql`

### Тело миграций

Все миграции должны иметь следущий вид:

```postgresql
SELECT assert_latest_migration(<last_migrations_id>);

-- код миграции

SELECT log_migration(<current_migrations_id>);
```

Где:
* **last_migrations_id** - номер предыдущей миграции,
* **current_migrations_id** - номер текущей миграции.

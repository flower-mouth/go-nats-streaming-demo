-- Удаление внешних ключей
ALTER TABLE IF EXISTS public.delivery DROP CONSTRAINT IF EXISTS delivery_order_uid_fkey;
ALTER TABLE IF EXISTS public.items DROP CONSTRAINT IF EXISTS items_order_uid_fkey;
ALTER TABLE IF EXISTS public.payment DROP CONSTRAINT IF EXISTS payment_order_uid_fkey;

-- Удаление таблиц
DROP TABLE IF EXISTS public.delivery;
DROP TABLE IF EXISTS public.items;
DROP TABLE IF EXISTS public.payment;
DROP TABLE IF EXISTS public.orders;

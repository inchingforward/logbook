-- Note that this creates an array with a length of one and inserts the original string.

alter table logbook_entry alter tags type text[] using array[tags];
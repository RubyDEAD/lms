// src/supabaseClient.js
import { createClient } from '@supabase/supabase-js'

const supabaseUrl = 'https://roifumfsdyhyegikikpd.supabase.co'
const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJvaWZ1bWZzZHloeWVnaWtpa3BkIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDA4NDQyNzIsImV4cCI6MjA1NjQyMDI3Mn0.h2KoDNHiVPCKnaEWjPADO0TdY4tmO1lGkssjuJDWxbs'

const bookserviceURL = 'https://hwkuzfsecehszlftxqpn.supabase.co'
const bookKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Imh3a3V6ZnNlY2Voc3psZnR4cXBuIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDA2MzcwNDgsImV4cCI6MjA1NjIxMzA0OH0.lCAde3mBhDNZg5kz9ffyF2xYJ2bb7BfOBjC1ZpbIQd4'
const storageURL = 'https://hwkuzfsecehszlftxqpn.supabase.co/storage/v1/object/public/test/'

export const supabase = createClient(supabaseUrl, supabaseAnonKey)

export const bookservice = createClient(bookserviceURL, bookKey)
export const bucketURL = storageURL
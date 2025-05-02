// src/supabaseClient.js
import { createClient } from '@supabase/supabase-js'

const supabaseUrl = 'https://roifumfsdyhyegikikpd.supabase.co'
const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InJvaWZ1bWZzZHloeWVnaWtpa3BkIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDA4NDQyNzIsImV4cCI6MjA1NjQyMDI3Mn0.h2KoDNHiVPCKnaEWjPADO0TdY4tmO1lGkssjuJDWxbs'

export const supabase = createClient(supabaseUrl, supabaseAnonKey)

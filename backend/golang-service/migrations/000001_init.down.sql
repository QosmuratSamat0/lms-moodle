-- Drop in reverse dependency order
DROP TABLE IF EXISTS sessions;

DROP TABLE IF EXISTS plagiarism_reports;

DROP TABLE IF EXISTS analytics_events;

DROP TABLE IF EXISTS schedule_events;

DROP TABLE IF EXISTS notifications;

DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_room_members;
DROP TABLE IF EXISTS chat_rooms;

DROP TABLE IF EXISTS attendance_marks;
DROP TABLE IF EXISTS attendance_sessions;

DROP TABLE IF EXISTS grades;

DROP TABLE IF EXISTS submissions;

DROP TABLE IF EXISTS assignments;

DROP TABLE IF EXISTS enrollments;

DROP TABLE IF EXISTS courses;

DROP TABLE IF EXISTS managers;
DROP TABLE IF EXISTS teachers;
DROP TABLE IF EXISTS students;

DROP TABLE IF EXISTS users;

-- Optional: extension stays (usually not dropped)
-- DROP EXTENSION IF EXISTS "uuid-ossp";

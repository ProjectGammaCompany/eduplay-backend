import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'https://api.hse-eduplay.ru';

const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXIxQG1haWwucnUiLCJleHAiOjE3NzY5NjU4MzgsImlkIjoiNjZmMzNmOWYtMDlhMC00NzA3LWFjYTktYWE3Nzg3NmIxMmFhIn0.foZkv9BhZ_UC3DDyd5Q2YIDFxSLb2lSQy0vWa3TFlhc';

export const options = {
    stages: [
        { duration: '30s', target: 10 },   // 10 пользователей
        { duration: '30s', target: 30 },   // 30 пользователей
        { duration: '30s', target: 50 },   // 50 пользователей
        { duration: '30s', target: 0 },    // остановка
    ],
};

export default function () {
    const res = http.get(`${BASE_URL}/events?page=1&maxOnPage=10`, {
        headers: {
            Authorization: `Bearer ${TOKEN}`,
        },
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(1);
}
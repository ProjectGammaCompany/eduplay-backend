import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'https://api.hse-eduplay.ru';

const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXI0QG1haWwucnUiLCJleHAiOjE3NzcyMjU3MjcsImlkIjoiZWE5YjYwOTItNDg5Mi00Y2E4LTk4ZTItNjY0M2I1NjIwZjg1In0.CtlDesf8V1vQ3Pcw0-AsJLmhtPs_bwufMRtXRv4bY1k';

export const options = {
    stages: [
        { duration: '30s', target: 10 },   // 10 пользователей
        { duration: '30s', target: 30 },   // 30 пользователей
        { duration: '30s', target: 50 },   // 50 пользователей
        { duration: '30s', target: 0 },    // остановка
    ],
};

export default function () {
    const res = http.get(`${BASE_URL}/event/648cc37b-bc2b-4308-aef0-426fc23d3b43/nextStage`, {
        headers: {
            Authorization: `Bearer ${TOKEN}`,
        },
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(1);
}
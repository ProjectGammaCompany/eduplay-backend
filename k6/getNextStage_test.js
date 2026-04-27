import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'https://api.hse-eduplay.ru';

const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXIxQG1haWwucnUiLCJleHAiOjE3NzcyNDcxNTYsImlkIjoiNjZmMzNmOWYtMDlhMC00NzA3LWFjYTktYWE3Nzg3NmIxMmFhIn0.kTnF6y8D61Jjkxh7E7jY3PuE91OfqTJNHEJlcZKpYp8';

export const options = {
    stages: [
        { duration: '1m', target: 10 },   
        { duration: '2m', target: 50 },   
        { duration: '30s', target: 100 },   
        { duration: '1m', target: 0 },    
    ],
};

export default function () {
    const res = http.get(`${BASE_URL}/event/5e4e87af-17f8-4308-b3c7-4787488c005f/nextStage`, {
        headers: {
            Authorization: `Bearer ${TOKEN}`,
        },
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(1);
}
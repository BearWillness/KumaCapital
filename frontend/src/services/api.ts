import axios from 'axios';

const apiClient = axios.create({
    baseURL: 'http://localhost:8080/atlas',
    headers: {
        'Content-Type': 'application/json',
    },
});

export const fetchUnemploymentRate = () => apiClient.get('/unemployment');
export const fetchInflationRate = () => apiClient.get('/inflation');
export const fetchInterestRate = () => apiClient.get('/interest_rate');
export const fetchGDPGrowth = () => apiClient.get('/gdp_growth');

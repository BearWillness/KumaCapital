import React, { useEffect, useState } from 'react';
import { fetchUnemploymentRate, fetchInflationRate, fetchInterestRate, fetchGDPGrowth } from '../services/api';

const EconomicDataDisplay = () => {
    const [unemploymentData, setUnemploymentData] = useState(null);
    const [inflationData, setInflationData] = useState(null);
    const [interestRateData, setInterestRateData] = useState(null);
    const [gdpGrowthData, setGdpGrowthData] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const unemployment = await fetchUnemploymentRate();
                setUnemploymentData(unemployment.data);

                const inflation = await fetchInflationRate();
                setInflationData(inflation.data);

                const interestRate = await fetchInterestRate();
                setInterestRateData(interestRate.data);

                const gdpGrowth = await fetchGDPGrowth();
                setGDPGrowthData(gdpGrowth.data);
            } catch (error) {
                console.error('Error fetching economic data:', error);
            }
        };

        fetchData();
    }, []);

    if (!unemploymentData || !inflationData || !interestRateData || !gdpGrowthData) {
        return <div>Loading...</div>;
    }

    return (
        <div>
            <h2>Economic Data</h2>
            <div>
                <h3>{unemploymentData.label}</h3>
                <p>Value: {unemploymentData.value}</p>
                <p>Risk: {unemploymentData.risk}%</p>
                <p>Recommendation: {unemploymentData.recommendation}</p>
            </div>
            <div>
                <h3>{inflationData.label}</h3>
                <p>Value: {inflationData.value}%</p>
                <p>Risk: {inflationData.risk}%</p>
                <p>Recommendation: {inflationData.recommendation}</p>
            </div>
            <div>
                <h3>{interestRateData.label}</h3>
                <p>Value: {interestRateData.value}%</p>
                <p>Risk: {interestRateData.risk}%</p>
                <p>Recommendation: {interestRateData.recommendation}</p>
            </div>
            <div>
                <h3>{gdpGrowthData.label}</h3>
                <p>Value: {gdpGrowthData.value}%</p>
                <p>Risk: {gdpGrowthData.risk}%</p>
                <p>Recommendation: {gdpGrowthData.recommendation}</p>
            </div>
        </div>
    );
};

export default EconomicDataDisplay;

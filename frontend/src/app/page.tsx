"use client";

import React, { useEffect, useState } from 'react';
import { fetchUnemploymentRate, fetchInflationRate, fetchInterestRate, fetchGDPGrowth } from '@/services/api';

export default function Home() {
  const [unemploymentData, setUnemploymentData] = useState(null);
  const [inflationData, setInflationData] = useState(null);
  const [interestRateData, setInterestRateData] = useState(null);
  const [gdpGrowthData, setGdpGrowthData] = useState(null); 

  useEffect(() => {
    const fetchData = async () => {
      try {
        const unemployment = await fetchUnemploymentRate();
        console.log('Unemployment Data:', unemployment.data);
        setUnemploymentData(unemployment.data);

        const inflation = await fetchInflationRate();
        console.log('Inflation Data:', inflation.data);
        setInflationData(inflation.data);

        const interestRate = await fetchInterestRate();
        console.log('Interest Rate Data:', interestRate.data);
        setInterestRateData(interestRate.data);

        const gdpGrowth = await fetchGDPGrowth();
        console.log('GDP Growth Data:', gdpGrowth.data);
        setGdpGrowthData(gdpGrowth.data); 
      } catch (error) {
        console.error('Error fetching economic data:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <h1 className="text-4xl font-bold mb-8">Economic Data Dashboard</h1>
      <div className="grid gap-8 lg:grid-cols-2">
        {unemploymentData && (
          <div className="border rounded-lg p-4 shadow-lg">
            <h2 className="text-2xl font-semibold">{unemploymentData.label}</h2>
            <p>Value: {unemploymentData.value}%</p>
            <p>Risk: {unemploymentData.risk}%</p>
            <p>Recommendation: {unemploymentData.recommendation}</p>
          </div>
        )}
        {inflationData && (
          <div className="border rounded-lg p-4 shadow-lg">
            <h2 className="text-2xl font-semibold">{inflationData.label}</h2>
            <p>Value: {inflationData.value}%</p>
            <p>Risk: {inflationData.risk}%</p>
            <p>Recommendation: {inflationData.recommendation}</p>
          </div>
        )}
        {interestRateData && (
          <div className="border rounded-lg p-4 shadow-lg">
            <h2 className="text-2xl font-semibold">{interestRateData.label}</h2>
            <p>Value: {interestRateData.value}%</p>
            <p>Risk: {interestRateData.risk}%</p>
            <p>Recommendation: {interestRateData.recommendation}</p>
          </div>
        )}
        {gdpGrowthData && (
          <div className="border rounded-lg p-4 shadow-lg">
            <h2 className="text-2xl font-semibold">{gdpGrowthData.label}</h2>
            <p>Value: {gdpGrowthData.value}%</p>
            <p>Risk: {gdpGrowthData.risk}%</p>
            <p>Recommendation: {gdpGrowthData.recommendation}</p>
          </div>
        )}
      </div>
    </main>
  );
}

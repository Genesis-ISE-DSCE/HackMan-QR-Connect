import React, { useState } from "react";
import { useParams } from "react-router-dom";
import "tailwindcss/tailwind.css";

interface Participant {
  teamName: string;
  name: string;
  email: string;
  phoneNum: string;
  breakfast: boolean;
  lunch: boolean;
  dinner: boolean;
  snack1: boolean;
  snack2: boolean;
}

const Dashboard: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  const [participant, setParticipant] = useState<Participant>({
    teamName: "trying..",
    name: "Aditya Agarwal",
    email: "adi790u@gmail.com",
    phoneNum: "9179822431",
    breakfast: true,
    lunch: true,
    dinner: false,
    snack1: true,
    snack2: false,
  });

  const toggleMealStatus = (meal: keyof Participant) => {
    setParticipant((prevState) => ({
      ...prevState,
      [meal]: !prevState[meal],
    }));
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-800 p-6">
      <div className=" text-gray-800 p-8 max-w-xl w-full">
        <h1 className="text-4xl md:text-6xl font-extralight mb-6 text-center tracking-wider bg-gradient-to-r from-purple-600 via-pink-600 to-red-600 bg-clip-text text-transparent p-4 animate-pulse">
          Participant Details
        </h1>
        <div className="space-y-4">
          <div className="flex flex-col mb-4">
            <h2 className="text-lg font-semibold text-neutral-200">
              Team Name
            </h2>
            <p className="text-gray-700 text-xl bg-gray-100 rounded-lg p-2 mt-1">
              {participant.teamName}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-lg font-semibold text-neutral-200">Name</h2>
            <p className="text-gray-700 text-xl bg-gray-100 rounded-lg p-2 mt-1">
              {participant.name}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-lg font-semibold text-neutral-200">Email</h2>
            <p className="text-gray-700 text-xl bg-gray-100 rounded-lg p-2 mt-1">
              {participant.email}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-lg font-semibold text-neutral-200">
              Phone Number
            </h2>
            <p className="text-gray-700 text-xl bg-gray-100 rounded-lg p-2 mt-1">
              {participant.phoneNum}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-lg font-semibold text-neutral-200">Meals</h2>
            <ul className="text-gray-700 space-y-2 mt-1">
              {[
                { meal: "breakfast", label: "Breakfast" },
                { meal: "lunch", label: "Lunch" },
                { meal: "dinner", label: "Dinner" },
                { meal: "snack1", label: "Snack 1" },
                { meal: "snack2", label: "Snack 2" },
              ].map(({ meal, label }) => (
                <li key={meal} className="flex items-center justify-between">
                  <label className="flex items-center space-x-3">
                    <input
                      type="checkbox"
                      //@ts-ignore
                      checked={participant[meal]}
                      onChange={() =>
                        toggleMealStatus(meal as keyof Participant)
                      }
                      className="form-checkbox h-5 w-5 text-neutral-200"
                    />
                    <span className="text-neutral-200">{label}</span>
                  </label>
                  <span
                    className={`px-2 py-1 text-sm rounded-full ${
                      //@ts-ignore
                      participant[meal]
                        ? "bg-green-100 text-green-800"
                        : "bg-red-100 text-red-800"
                    }`}
                  >
                    {/* @ts-ignore */}
                    {participant[meal] ? "Yes" : "No"}
                  </span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;

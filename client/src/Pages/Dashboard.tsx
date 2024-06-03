import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import "tailwindcss/tailwind.css";
import connector from "../api";

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
  const navigate = useNavigate();

  const [participant, setParticipant] = useState<Participant>({
    teamName: "",
    name: "",
    email: "",
    phoneNum: "",
    breakfast: false,
    lunch: false,
    dinner: false,
    snack1: false,
    snack2: false,
  });

  useEffect(() => {
    const fetchData = async () => {
      const token = localStorage.getItem("token");

      if (!token) {
        navigate("/");
      }
      const response = await connector.get(`/user/details/${id}`);
      const data = response.data.participant;

      setParticipant({
        teamName: data.TeamName,
        name: data.Name,
        email: data.Email,
        phoneNum: data.PhoneNum,
        breakfast: data.Breakfast,
        lunch: data.Lunch,
        dinner: data.Dinner,
        snack1: data.Snack1,
        snack2: data.Snack2,
      });
    };
    fetchData();
  }, [id]);

  const toggleMealStatus = async (meal: keyof Participant) => {
    const mealData = {
      Meal: meal,
    };
    await connector.post(`/user/update/${id}`, mealData);

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
            <h2 className="text-2xl font-semibold text-neutral-200 tracking-wider">
              Team Name :
            </h2>
            <p className="text-gray-300 text-lg rounded-lg mt-2 font-extralight font-mono overflow-auto">
              {participant.teamName}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-2xl font-semibold text-neutral-200 tracking-wider">
              Name :
            </h2>
            <p className="text-gray-300 text-lg rounded-lg mt-2 font-extralight font-mono overflow-auto">
              {participant.name}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-2xl font-semibold text-neutral-200 tracking-wider">
              Email :
            </h2>
            <p className="text-gray-300 text-lg rounded-lg mt-2 font-extralight font-mono overflow-auto">
              {participant.email}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-2xl font-semibold text-neutral-200 tracking-wider">
              Phone Number :
            </h2>
            <p className="text-gray-300 text-lg rounded-lg mt-2 font-extralight font-mono">
              {participant.phoneNum}
            </p>
          </div>
          <div className="flex flex-col mb-4">
            <h2 className="text-2xl font-semibold text-neutral-200 tracking-wider">
              Meals :
            </h2>
            <ul className="text-gray-700 space-y-2 mt-2">
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

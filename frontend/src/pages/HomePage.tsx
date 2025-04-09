import React from "react";
import Layout from "../components/Layout";
import Wrapper from "../components/Wrapper";

const HomePage: React.FC = () => {
  return (
    <Layout>
      <Wrapper title="Home Page">
        <p className="mb-4">Welcome to the application!</p>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
          <div className="bg-blue-50/80 p-4 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-blue-700 mb-2">
              Feature 1
            </h3>
            <p>
              Description of an amazing feature that our application provides.
            </p>
          </div>
          <div className="bg-green-50/80 p-4 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-green-700 mb-2">
              Feature 2
            </h3>
            <p>
              Another incredible feature that makes our application special.
            </p>
          </div>
        </div>
      </Wrapper>
    </Layout>
  );
};

export default HomePage;

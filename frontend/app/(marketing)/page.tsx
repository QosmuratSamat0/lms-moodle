"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import {
  ArrowRight,
  BookOpen,
  Users,
  Calendar,
  MessageSquare,
  Award,
  BarChart3,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useUIStore } from "@/store/ui-store";

const features = [
  {
    icon: BookOpen,
    title: "Course Management",
    description:
      "Create and manage courses with lectures, assignments, and quizzes.",
  },
  {
    icon: Users,
    title: "Collaborative Learning",
    description: "Connect students and teachers in an interactive environment.",
  },
  {
    icon: Calendar,
    title: "Schedule & Deadlines",
    description: "Keep track of classes, exams, and assignment due dates.",
  },
  {
    icon: MessageSquare,
    title: "Real-time Chat",
    description: "Communicate instantly with peers and instructors.",
  },
  {
    icon: Award,
    title: "Grading & Feedback",
    description: "Comprehensive grading system with detailed feedback.",
  },
  {
    icon: BarChart3,
    title: "Progress Analytics",
    description: "Track learning progress with insightful analytics.",
  },
];

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
    },
  },
};

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.5,
    },
  },
};

export default function HomePage() {
  const { openLoginModal } = useUIStore();

  return (
    <>
      {/* Hero Section */}
      <section className="relative overflow-hidden py-12 sm:py-16 md:py-24 lg:py-32">
        <div className="absolute inset-0 -z-10 bg-[linear-gradient(to_right,#f0f0f0_1px,transparent_1px),linear-gradient(to_bottom,#f0f0f0_1px,transparent_1px)] dark:bg-[linear-gradient(to_right,#1a1a2e_1px,transparent_1px),linear-gradient(to_bottom,#1a1a2e_1px,transparent_1px)] bg-size-[4rem_4rem] mask-[radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_110%)]" />

        <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="mx-auto max-w-3xl text-center"
          >
            <motion.div
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 0.2, duration: 0.5 }}
              className="mb-4 sm:mb-6 inline-flex items-center rounded-full border bg-background px-3 sm:px-4 py-1 sm:py-1.5 text-xs sm:text-sm"
            >
              <span className="mr-1.5 sm:mr-2">🎓</span>
              <span>Welcome to the future of learning</span>
            </motion.div>

            <h1 className="text-3xl sm:text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-bold tracking-tight">
              Learn without{" "}
              <span className="bg-linear-to-r from-primary to-blue-600 bg-clip-text text-transparent">
                limits
              </span>
            </h1>

            <p className="mt-4 sm:mt-6 text-base sm:text-lg text-muted-foreground md:text-xl px-2 sm:px-0">
              A modern learning management system designed for students and
              teachers. Create courses, track progress, and collaborate in
              real-time.
            </p>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4, duration: 0.5 }}
              className="mt-6 sm:mt-8 md:mt-10 flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center px-4 sm:px-0"
            >
              <Button
                size="lg"
                onClick={openLoginModal}
                className="h-12 sm:h-11 text-base sm:text-sm touch-manipulation"
              >
                Login to Get Started
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </motion.div>
          </motion.div>

          {/* Hero Image - Dashboard Preview */}
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6, duration: 0.8 }}
            className="mx-auto mt-10 sm:mt-12 md:mt-16 max-w-5xl px-2 sm:px-0"
          >
            <div className="relative rounded-lg sm:rounded-xl border bg-linear-to-b from-muted/50 to-muted p-1.5 sm:p-2 shadow-xl sm:shadow-2xl overflow-hidden">
              <div className="rounded-md sm:rounded-lg overflow-hidden">
                <Image
                  src="/images/dashboard-preview.png"
                  alt="LMS Dashboard Preview"
                  width={1920}
                  height={1080}
                  className="w-full h-auto"
                  priority
                />
              </div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Features Section */}
      <section
        id="features"
        className="py-12 sm:py-16 md:py-24 lg:py-32 bg-muted/30"
      >
        <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="mx-auto max-w-2xl text-center"
          >
            <h2 className="text-2xl sm:text-3xl font-bold tracking-tight md:text-4xl">
              Everything you need to succeed
            </h2>
            <p className="mt-3 sm:mt-4 text-base sm:text-lg text-muted-foreground">
              Powerful features designed to enhance the learning experience
            </p>
          </motion.div>

          <motion.div
            variants={containerVariants}
            initial="hidden"
            whileInView="visible"
            viewport={{ once: true }}
            className="mx-auto mt-10 sm:mt-12 md:mt-16 grid gap-4 sm:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3"
          >
            {features.map((feature) => (
              <motion.div key={feature.title} variants={itemVariants}>
                <Card className="h-full transition-shadow hover:shadow-lg">
                  <CardHeader className="p-4 sm:p-6">
                    <div className="mb-2 inline-flex h-10 w-10 sm:h-12 sm:w-12 items-center justify-center rounded-lg bg-primary/10">
                      <feature.icon className="h-5 w-5 sm:h-6 sm:w-6 text-primary" />
                    </div>
                    <CardTitle className="text-base sm:text-lg">
                      {feature.title}
                    </CardTitle>
                    <CardDescription className="text-sm">
                      {feature.description}
                    </CardDescription>
                  </CardHeader>
                </Card>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-12 sm:py-16 md:py-24 lg:py-32">
        <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            whileInView={{ opacity: 1, scale: 1 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="mx-auto max-w-4xl rounded-xl sm:rounded-2xl bg-linear-to-r from-primary/10 via-primary/5 to-primary/10 p-6 sm:p-8 md:p-12 text-center border"
          >
            <h2 className="text-2xl sm:text-3xl font-bold tracking-tight md:text-4xl">
              Ready to transform your learning?
            </h2>
            <p className="mt-3 sm:mt-4 text-base sm:text-lg text-muted-foreground">
              Join thousands of students and teachers already using our
              platform.
            </p>
            <div className="mt-6 sm:mt-8 flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center">
              <Button
                size="lg"
                onClick={openLoginModal}
                className="h-12 sm:h-11 text-base sm:text-sm touch-manipulation"
              >
                Login to Your Account
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </div>
          </motion.div>
        </div>
      </section>
    </>
  );
}

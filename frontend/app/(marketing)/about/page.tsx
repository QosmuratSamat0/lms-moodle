"use client";

import { motion } from "framer-motion";
import { Target, Users, Lightbulb, Award } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const values = [
  {
    icon: Target,
    title: "Our Mission",
    description:
      "To democratize education by providing a powerful, accessible learning platform for everyone.",
  },
  {
    icon: Users,
    title: "Community First",
    description:
      "Building a vibrant community of learners and educators who support each other.",
  },
  {
    icon: Lightbulb,
    title: "Innovation",
    description:
      "Continuously improving our platform with cutting-edge technology and pedagogical insights.",
  },
  {
    icon: Award,
    title: "Excellence",
    description:
      "Striving for excellence in every feature we build and every experience we create.",
  },
];

const team = [
  {
    name: "Dr. Sarah Chen",
    role: "CEO & Founder",
    bio: "Former professor with 15 years of experience in educational technology.",
  },
  {
    name: "Michael Rodriguez",
    role: "CTO",
    bio: "Tech veteran who has built scalable platforms at leading EdTech companies.",
  },
  {
    name: "Emily Watson",
    role: "Head of Product",
    bio: "Product leader passionate about creating intuitive learning experiences.",
  },
  {
    name: "David Kim",
    role: "Head of Engineering",
    bio: "Engineering leader focused on building reliable and performant systems.",
  },
];

export default function AboutPage() {
  return (
    <>
      {/* Hero Section */}
      <section className="py-20 md:py-32">
        <div className="container">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
            className="mx-auto max-w-3xl text-center"
          >
            <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
              About Us
            </h1>
            <p className="mt-6 text-lg text-muted-foreground">
              We&apos;re on a mission to make quality education accessible to
              everyone, everywhere. Our platform connects students and teachers
              in meaningful ways.
            </p>
          </motion.div>
        </div>
      </section>

      {/* Values Section */}
      <section className="py-20 bg-muted/30">
        <div className="container">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="mx-auto max-w-2xl text-center mb-12"
          >
            <h2 className="text-3xl font-bold tracking-tight">Our Values</h2>
            <p className="mt-4 text-muted-foreground">
              The principles that guide everything we do
            </p>
          </motion.div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
            {values.map((value, index) => (
              <motion.div
                key={value.title}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
              >
                <Card className="h-full text-center">
                  <CardHeader>
                    <div className="mx-auto mb-2 inline-flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                      <value.icon className="h-6 w-6 text-primary" />
                    </div>
                    <CardTitle className="text-lg">{value.title}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground">
                      {value.description}
                    </p>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Team Section */}
      <section className="py-20">
        <div className="container">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="mx-auto max-w-2xl text-center mb-12"
          >
            <h2 className="text-3xl font-bold tracking-tight">Our Team</h2>
            <p className="mt-4 text-muted-foreground">
              Meet the people behind our platform
            </p>
          </motion.div>

          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
            {team.map((member, index) => (
              <motion.div
                key={member.name}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
              >
                <Card className="h-full">
                  <CardHeader className="text-center">
                    <div className="mx-auto mb-4 h-24 w-24 rounded-full bg-muted flex items-center justify-center">
                      <span className="text-3xl">👤</span>
                    </div>
                    <CardTitle className="text-lg">{member.name}</CardTitle>
                    <p className="text-sm text-primary">{member.role}</p>
                  </CardHeader>
                  <CardContent>
                    <p className="text-sm text-muted-foreground text-center">
                      {member.bio}
                    </p>
                  </CardContent>
                </Card>
              </motion.div>
            ))}
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-20 bg-muted/30">
        <div className="container">
          <div className="grid gap-8 md:grid-cols-4 text-center">
            {[
              { label: "Students", value: "50,000+" },
              { label: "Teachers", value: "2,000+" },
              { label: "Courses", value: "5,000+" },
              { label: "Countries", value: "120+" },
            ].map((stat, index) => (
              <motion.div
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: index * 0.1 }}
              >
                <div className="text-4xl font-bold text-primary">
                  {stat.value}
                </div>
                <div className="mt-2 text-muted-foreground">{stat.label}</div>
              </motion.div>
            ))}
          </div>
        </div>
      </section>
    </>
  );
}

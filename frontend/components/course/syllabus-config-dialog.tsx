"use client";

import { useState, useEffect } from "react";
import { Plus, Trash2, Info } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Alert, AlertDescription } from "@/components/ui/alert";

interface SyllabusComponent {
  name: string;
  count: number;
  weightEach: number;
}

interface SyllabusConfig {
  registerMidterm: {
    weight: number;
    components: SyllabusComponent[];
  };
  registerEndterm: {
    weight: number;
    components: SyllabusComponent[];
  };
  final: {
    weight: number;
  };
}

interface SyllabusConfigDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  config: SyllabusConfig;
  onSave: (config: SyllabusConfig) => void;
}

const componentTypes = [
  "Assignment",
  "Quiz",
  "Lab Work",
  "Project",
  "Presentation",
  "Homework",
];

export function SyllabusConfigDialog({
  open,
  onOpenChange,
  config,
  onSave,
}: SyllabusConfigDialogProps) {
  const [localConfig, setLocalConfig] = useState<SyllabusConfig>(config);

  useEffect(() => {
    setLocalConfig(config);
  }, [config]);

  const updateWeight = (
    section: "registerMidterm" | "registerEndterm" | "final",
    value: number,
  ) => {
    setLocalConfig((prev) => ({
      ...prev,
      [section]: {
        ...prev[section],
        weight: value,
      },
    }));
  };

  const addComponent = (section: "registerMidterm" | "registerEndterm") => {
    setLocalConfig((prev) => ({
      ...prev,
      [section]: {
        ...prev[section],
        components: [
          ...prev[section].components,
          { name: "Assignment", count: 1, weightEach: 0 },
        ],
      },
    }));
  };

  const removeComponent = (
    section: "registerMidterm" | "registerEndterm",
    index: number,
  ) => {
    setLocalConfig((prev) => ({
      ...prev,
      [section]: {
        ...prev[section],
        components: prev[section].components.filter((_, i) => i !== index),
      },
    }));
  };

  const updateComponent = (
    section: "registerMidterm" | "registerEndterm",
    index: number,
    field: keyof SyllabusComponent,
    value: string | number,
  ) => {
    setLocalConfig((prev) => ({
      ...prev,
      [section]: {
        ...prev[section],
        components: prev[section].components.map((comp, i) =>
          i === index ? { ...comp, [field]: value } : comp,
        ),
      },
    }));
  };

  // Auto-calculate weightEach when count changes
  const autoCalculateWeights = (
    section: "registerMidterm" | "registerEndterm",
  ) => {
    const sectionData = localConfig[section];
    const totalComponents = sectionData.components.reduce(
      (sum, c) => sum + c.count,
      0,
    );
    if (totalComponents > 0) {
      const calculatedWeight = Math.round(100 / totalComponents);
      setLocalConfig((prev) => ({
        ...prev,
        [section]: {
          ...prev[section],
          components: prev[section].components.map((comp) => ({
            ...comp,
            weightEach: calculatedWeight,
          })),
        },
      }));
    }
  };

  const totalWeight =
    localConfig.registerMidterm.weight +
    localConfig.registerEndterm.weight +
    localConfig.final.weight;
  const isValid = totalWeight === 100;

  const handleSave = () => {
    if (isValid) {
      onSave(localConfig);
      onOpenChange(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Configure Syllabus Grading</DialogTitle>
          <DialogDescription>
            Set up the grading structure for this course based on your syllabus
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* Weight Summary */}
          <Alert variant={isValid ? "default" : "destructive"}>
            <Info className="h-4 w-4" />
            <AlertDescription>
              Total weight: <strong>{totalWeight}%</strong>
              {!isValid && " (must equal 100%)"}
            </AlertDescription>
          </Alert>

          {/* Main Weights */}
          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-2">
              <Label>Register Midterm %</Label>
              <Input
                type="number"
                min={0}
                max={100}
                value={localConfig.registerMidterm.weight}
                onChange={(e) =>
                  updateWeight("registerMidterm", parseInt(e.target.value) || 0)
                }
              />
            </div>
            <div className="space-y-2">
              <Label>Register Endterm %</Label>
              <Input
                type="number"
                min={0}
                max={100}
                value={localConfig.registerEndterm.weight}
                onChange={(e) =>
                  updateWeight("registerEndterm", parseInt(e.target.value) || 0)
                }
              />
            </div>
            <div className="space-y-2">
              <Label>Final Exam %</Label>
              <Input
                type="number"
                min={0}
                max={100}
                value={localConfig.final.weight}
                onChange={(e) =>
                  updateWeight("final", parseInt(e.target.value) || 0)
                }
              />
            </div>
          </div>

          {/* Components Configuration */}
          <Accordion
            type="multiple"
            defaultValue={["midterm", "endterm"]}
            className="w-full"
          >
            {/* Register Midterm Components */}
            <AccordionItem value="midterm">
              <AccordionTrigger className="text-blue-700 dark:text-blue-400">
                Register Midterm Components (
                {localConfig.registerMidterm.weight}%)
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-3 pt-2">
                  {localConfig.registerMidterm.components.map((comp, idx) => (
                    <div
                      key={idx}
                      className="flex items-end gap-2 p-3 rounded border bg-blue-50/50 dark:bg-blue-900/10"
                    >
                      <div className="flex-1 space-y-1">
                        <Label className="text-xs">Type</Label>
                        <Select
                          value={comp.name}
                          onValueChange={(v) =>
                            updateComponent("registerMidterm", idx, "name", v)
                          }
                        >
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {componentTypes.map((type) => (
                              <SelectItem key={type} value={type}>
                                {type}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="w-20 space-y-1">
                        <Label className="text-xs">Count</Label>
                        <Input
                          type="number"
                          min={1}
                          value={comp.count}
                          onChange={(e) =>
                            updateComponent(
                              "registerMidterm",
                              idx,
                              "count",
                              parseInt(e.target.value) || 1,
                            )
                          }
                        />
                      </div>
                      <div className="w-24 space-y-1">
                        <Label className="text-xs">Weight Each %</Label>
                        <Input
                          type="number"
                          min={0}
                          max={100}
                          value={comp.weightEach}
                          onChange={(e) =>
                            updateComponent(
                              "registerMidterm",
                              idx,
                              "weightEach",
                              parseInt(e.target.value) || 0,
                            )
                          }
                        />
                      </div>
                      <Button
                        size="icon"
                        variant="ghost"
                        className="text-red-500 hover:text-red-700"
                        onClick={() => removeComponent("registerMidterm", idx)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => addComponent("registerMidterm")}
                    >
                      <Plus className="h-4 w-4 mr-1" />
                      Add Component
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => autoCalculateWeights("registerMidterm")}
                    >
                      Auto-calculate weights
                    </Button>
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>

            {/* Register Endterm Components */}
            <AccordionItem value="endterm">
              <AccordionTrigger className="text-purple-700 dark:text-purple-400">
                Register Endterm Components (
                {localConfig.registerEndterm.weight}%)
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-3 pt-2">
                  {localConfig.registerEndterm.components.map((comp, idx) => (
                    <div
                      key={idx}
                      className="flex items-end gap-2 p-3 rounded border bg-purple-50/50 dark:bg-purple-900/10"
                    >
                      <div className="flex-1 space-y-1">
                        <Label className="text-xs">Type</Label>
                        <Select
                          value={comp.name}
                          onValueChange={(v) =>
                            updateComponent("registerEndterm", idx, "name", v)
                          }
                        >
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            {componentTypes.map((type) => (
                              <SelectItem key={type} value={type}>
                                {type}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="w-20 space-y-1">
                        <Label className="text-xs">Count</Label>
                        <Input
                          type="number"
                          min={1}
                          value={comp.count}
                          onChange={(e) =>
                            updateComponent(
                              "registerEndterm",
                              idx,
                              "count",
                              parseInt(e.target.value) || 1,
                            )
                          }
                        />
                      </div>
                      <div className="w-24 space-y-1">
                        <Label className="text-xs">Weight Each %</Label>
                        <Input
                          type="number"
                          min={0}
                          max={100}
                          value={comp.weightEach}
                          onChange={(e) =>
                            updateComponent(
                              "registerEndterm",
                              idx,
                              "weightEach",
                              parseInt(e.target.value) || 0,
                            )
                          }
                        />
                      </div>
                      <Button
                        size="icon"
                        variant="ghost"
                        className="text-red-500 hover:text-red-700"
                        onClick={() => removeComponent("registerEndterm", idx)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => addComponent("registerEndterm")}
                    >
                      <Plus className="h-4 w-4 mr-1" />
                      Add Component
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => autoCalculateWeights("registerEndterm")}
                    >
                      Auto-calculate weights
                    </Button>
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>
          </Accordion>

          {/* Example Calculation */}
          <div className="p-4 rounded-lg bg-muted text-sm">
            <p className="font-medium mb-2">Example Calculation:</p>
            <p className="text-muted-foreground">
              If a student scores 85% on Assignment 1 and 75% on Assignment 2 in
              Register Midterm:
            </p>
            <p className="text-muted-foreground">
              Midterm Average: (85 + 75) / 2 = 80%
            </p>
            <p className="text-muted-foreground">
              Contribution to total: 80% × {localConfig.registerMidterm.weight}%
              = {((80 * localConfig.registerMidterm.weight) / 100).toFixed(1)}%
            </p>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={!isValid}>
            Save Configuration
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

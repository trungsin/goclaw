import { useTranslation } from "react-i18next";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export function defaultAcpArgsForBinary(binary: string, currentArgs: string): string {
  if (currentArgs.trim()) return currentArgs;
  const base = binary.trim().split(/[\\/]/).pop()?.replace(/\.exe$/i, "") ?? "";
  if (base === "grok") return "agent stdio";
  return currentArgs;
}

interface ACPSectionProps {
  binary: string;
  onBinaryChange: (v: string) => void;
  args: string;
  onArgsChange: (v: string) => void;
  idleTTL: string;
  onIdleTTLChange: (v: string) => void;
  permMode: string;
  onPermModeChange: (v: string) => void;
  workDir: string;
  onWorkDirChange: (v: string) => void;
}

export function ACPSection({
  binary, onBinaryChange,
  args, onArgsChange,
  idleTTL, onIdleTTLChange,
  permMode, onPermModeChange,
  workDir, onWorkDirChange,
}: ACPSectionProps) {
  const { t } = useTranslation("providers");

  return (
    <div className="space-y-4">
      <p className="text-sm text-muted-foreground">{t("acp.description")}</p>

      <div className="space-y-2">
        <Label htmlFor="acpBinary">{t("acp.binary")}</Label>
        <Input
          id="acpBinary"
          value={binary}
          onChange={(e) => {
            const next = e.target.value;
            onBinaryChange(next);
            const filled = defaultAcpArgsForBinary(next, args);
            if (filled !== args) onArgsChange(filled);
          }}
          placeholder={t("acp.binaryPlaceholder")}
          className="text-base md:text-sm"
        />
        <p className="text-xs text-muted-foreground">{t("acp.binaryHint")}</p>
      </div>

      <div className="space-y-2">
        <Label htmlFor="acpArgs">{t("acp.args")}</Label>
        <Input
          id="acpArgs"
          value={args}
          onChange={(e) => onArgsChange(e.target.value)}
          placeholder={binary.trim().split(/[\\/]/).pop()?.replace(/\.exe$/i, "") === "grok" ? "agent stdio" : t("acp.argsPlaceholder")}
          className="text-base md:text-sm"
        />
        <p className="text-xs text-muted-foreground">{t("acp.argsHint")}</p>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="acpIdleTTL">{t("acp.idleTTL")}</Label>
          <Input
            id="acpIdleTTL"
            value={idleTTL}
            onChange={(e) => onIdleTTLChange(e.target.value)}
            placeholder={t("acp.idleTTLPlaceholder")}
            className="text-base md:text-sm"
          />
          <p className="text-xs text-muted-foreground">{t("acp.idleTTLHint")}</p>
        </div>

        <div className="space-y-2">
          <Label>{t("acp.permMode")}</Label>
          <Select value={permMode} onValueChange={onPermModeChange}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="approve-all">{t("acp.permModeApproveAll")}</SelectItem>
              <SelectItem value="approve-reads">{t("acp.permModeApproveReads")}</SelectItem>
              <SelectItem value="deny-all">{t("acp.permModeDenyAll")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="acpWorkDir">{t("acp.workDir")}</Label>
        <Input
          id="acpWorkDir"
          value={workDir}
          onChange={(e) => onWorkDirChange(e.target.value)}
          placeholder={t("acp.workDirPlaceholder")}
          className="text-base md:text-sm"
        />
        <p className="text-xs text-muted-foreground">{t("acp.workDirHint")}</p>
      </div>
    </div>
  );
}

package com.sourcegraph.website;

import com.intellij.openapi.actionSystem.AnAction;
import com.intellij.openapi.actionSystem.AnActionEvent;
import com.intellij.openapi.actionSystem.CommonDataKeys;
import com.intellij.openapi.application.ApplicationInfo;
import com.intellij.openapi.diagnostic.Logger;
import com.intellij.openapi.editor.Document;
import com.intellij.openapi.editor.Editor;
import com.intellij.openapi.editor.SelectionModel;
import com.intellij.openapi.fileEditor.FileDocumentManager;
import com.intellij.openapi.fileEditor.FileEditorManager;
import com.intellij.openapi.project.Project;
import com.intellij.openapi.vfs.VirtualFile;
import com.intellij.psi.PsiDocumentManager;
import com.intellij.psi.PsiElement;
import com.intellij.psi.PsiFile;
import com.sourcegraph.config.ConfigUtil;
import com.sourcegraph.git.GitUtil;
import com.sourcegraph.git.RepoInfo;
import org.jetbrains.annotations.Nullable;

import java.awt.*;
import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

public abstract class SearchActionBase extends AnAction {
    public void actionPerformedMode(AnActionEvent e, String mode) {
        Logger logger = Logger.getInstance(this.getClass());

        // Get project, editor, document, file, and position information.
        final Project project = e.getProject();
        if (project == null) {
            return;
        }
        Editor editor = FileEditorManager.getInstance(project).getSelectedTextEditor();
        if (editor == null) {
            return;
        }
        Document currentDoc = editor.getDocument();
        VirtualFile currentFile = FileDocumentManager.getInstance().getFile(currentDoc);
        if (currentFile == null) {
            return;
        }
        SelectionModel sel = editor.getSelectionModel();

        // Get repo information.
        RepoInfo repoInfo = GitUtil.getRepoInfo(currentFile.getParent().getPath(), project);

        String q = sel.getSelectedText();
        if (q == null || q.equals("")) {
            // If no selection check if identifier under caret and use it for search.
            PsiFile psiFile = PsiDocumentManager.getInstance(project).getPsiFile(currentDoc);
            if (psiFile == null) {
                return;
            }

            PsiElement psiElement = psiFile.findElementAt(editor.getCaretModel().getOffset());
            if (psiElement == null) {
                return;
            }

            q = psiElement.getText();

            if (q == null || q.isBlank() || q.length() == 1) {
                return; // nothing to query
            }
        }

        // Build the URL that we will open.
        String uri;
        String productName = ApplicationInfo.getInstance().getVersionName();
        String productVersion = ApplicationInfo.getInstance().getFullVersion();

        uri = ConfigUtil.getSourcegraphUrl(project) + "-/editor"
            + "?editor=" + URLEncoder.encode("JetBrains", StandardCharsets.UTF_8)
            + "&version=" + URLEncoder.encode(ConfigUtil.getVersion(), StandardCharsets.UTF_8)
            + "&utm_product_name=" + URLEncoder.encode(productName, StandardCharsets.UTF_8)
            + "&utm_product_version=" + URLEncoder.encode(productVersion, StandardCharsets.UTF_8)
            + "&search=" + URLEncoder.encode(q, StandardCharsets.UTF_8);

        if (mode.equals("search.repository")) {
            uri += "&search_remote_url=" + URLEncoder.encode(repoInfo.remoteUrl, StandardCharsets.UTF_8)
                + "&search_branch=" + URLEncoder.encode(repoInfo.branchName, StandardCharsets.UTF_8);
        }

        // Open the URL in the browser.
        try {
            Desktop.getDesktop().browse(URI.create(uri));
        } catch (IOException err) {
            logger.debug("failed to open browser");
            err.printStackTrace();
        }
    }

    @Override
    public void update(AnActionEvent e) {
        String selectedText = getSelectedText(e);
        e.getPresentation().setEnabled(selectedText != null && selectedText.length() > 0);
    }

    @Nullable
    private String getSelectedText(AnActionEvent e) {
        final Project project = e.getProject();
        if (project == null) {
            return null;
        }
        Editor editor = FileEditorManager.getInstance(project).getSelectedTextEditor();
        if (editor == null) {
            return null;
        }
        Document currentDoc = editor.getDocument();
        VirtualFile currentFile = FileDocumentManager.getInstance().getFile(currentDoc);
        if (currentFile == null) {
            return null;
        }
        SelectionModel sel = editor.getSelectionModel();

        String q = sel.getSelectedText();
        if (q == null || q.equals("")) {
            // If no selection check if identifier under caret and use it for search.
            PsiFile psiFile = PsiDocumentManager.getInstance(project).getPsiFile(currentDoc);
            if (psiFile == null) {
                return null;
            }

            PsiElement psiElement = psiFile.findElementAt(editor.getCaretModel().getOffset());
            if (psiElement == null) {
                return null;
            }

            q = psiElement.getText();

            if (q.isBlank() || q.length() == 1) {
                return null;
            }
        }

        return q;
    }
}
